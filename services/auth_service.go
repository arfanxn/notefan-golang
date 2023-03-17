package services

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/clarketm/json"
	"github.com/google/uuid"
	media_collnames "github.com/notefan-golang/enums/media/collection_names"
	media_disks "github.com/notefan-golang/enums/media/disks"
	token_types "github.com/notefan-golang/enums/token/types"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/chanh"
	"github.com/notefan-golang/helpers/jwth"
	"github.com/notefan-golang/helpers/mailh"
	"github.com/notefan-golang/helpers/numberh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/auth_reqs"
	authReqs "github.com/notefan-golang/models/requests/auth_reqs"
	"github.com/notefan-golang/models/requests/file_reqs"
	authRess "github.com/notefan-golang/models/responses/auth_ress"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/user_ress"
	"github.com/notefan-golang/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository  *repositories.UserRepository
	tokenRepository *repositories.TokenRepository
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.Mutex
}

func NewAuthService(
	userRepository *repositories.UserRepository,
	tokenRepository *repositories.TokenRepository,
	mediaRepository *repositories.MediaRepository) *AuthService {
	return &AuthService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		mediaRepository: mediaRepository,
		waitGroup:       new(sync.WaitGroup),
		mutex:           sync.Mutex{},
	}
}

func (service *AuthService) Login(ctx context.Context, data authReqs.Login) (loginRes authRess.Login, err error) {
	user, err := service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil || user.Id == uuid.Nil { // if err / user not found then return exception HTTPAuthLoginFailed
		err = exceptions.HTTPAuthLoginFailed
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil { // err not nil == password doesnt match, return exception HTTPAuthLoginFailed
		err = exceptions.HTTPAuthLoginFailed
		return
	}

	// Prepare claims (payload) for the JWT
	claims := map[string]any{}
	claims["authorized"] = true
	claims["id"] = user.Id.String()
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // expires in N minutes

	// get app key (signing key)
	signature := os.Getenv("APP_KEY")

	// Encode/Generate JWT token
	token, err := jwth.Encode(signature, claims)
	if err != nil {
		return
	}

	loginRes = authRess.Login{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}
	return loginRes, nil
}

// Register registers the given user
func (service *AuthService) Register(ctx context.Context, data authReqs.Register) (userRes user_ress.User, err error) {
	// Hash the user password
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	var (
		userEty  entities.User
		mediaEty entities.Media
		errChan  = chanh.Make[error](nil, 1)
	)

	service.waitGroup.Add(2)

	go func() { // goroutine for creating user
		defer service.waitGroup.Done()

		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		userEty, errChanVal := service.userRepository.FindByEmail(ctx, data.Email)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if userEty.Email == data.Email {
			bytes, errChanVal := json.Marshal(map[string]string{
				"email": "email already exists",
			})
			if errChanVal != nil {
				chanh.ReplaceVal(errChan, errChanVal)
				return
			}
			errChan <- exceptions.NewHTTPError(
				http.StatusUnprocessableEntity,
				errors.New(string(bytes)),
			)
			return
		}

		service.mutex.Lock()
		defer service.mutex.Unlock()
		// parse request to entity
		userEty = entities.User{
			Name:     data.Name,
			Email:    data.Email,
			Password: string(password),
		}

		// Save user into Database
		_, errChanVal = service.userRepository.Create(ctx, &userEty)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
	}()

	go func() { // goroutine for creating user's avatar
		defer service.waitGroup.Done()

		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		// the user default avatar file path
		defaultAvatarFilePath := "./public/placeholders/images/default-avatar.jpg"

		service.mutex.Lock()
		defer service.mutex.Unlock()

		// open default user avatar file
		defaultAvatarBytes, errChanVal := os.ReadFile(defaultAvatarFilePath)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}

		// Prepare media entity for user avatar
		mediaEty = entities.Media{
			ModelType:      reflecth.GetTypeName(userEty),
			ModelId:        userEty.Id,
			CollectionName: media_collnames.Avatar,
			Disk:           media_disks.Public,
			FileName:       filepath.Base(defaultAvatarFilePath),
			File:           file_reqs.NewFromBytes(defaultAvatarBytes),
		}

		// Save media into Database
		_, errChanVal = service.mediaRepository.Create(ctx, &mediaEty)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
	}()

	service.waitGroup.Wait()

	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}

	userRes = user_ress.FillFromEntity(userEty)
	userRes.Avatar = media_ress.FillFromEntity(mediaEty)

	// Return the created user and nil
	return userRes, nil
}

// ForgotPassword sends reset password otp to the given email from request
func (service *AuthService) ForgotPassword(ctx context.Context, data auth_reqs.ForgotPassword) (err error) {
	var (
		userEty  entities.User
		tokenEty entities.Token
	)

	// User relateds
	userEty, err = service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		return
	}
	if userEty.Id == uuid.Nil { // return not found if not found
		err = exceptions.HTTPNotFound
		return
	}

	// Token relateds
	tokenEty, err = service.tokenRepository.FindByTokenableAndType(
		ctx,
		reflecth.GetTypeName(userEty),
		userEty.Id.String(),
		token_types.ResetPassword,
	)
	if err != nil {
		return
	}
	// if token not found or has been used or expired then delete and create a new one
	if (tokenEty.Id == uuid.Nil) || (tokenEty.UsedAt.Valid) || tokenEty.IsExpired() {
		_, err = service.tokenRepository.DeleteByIds(ctx, tokenEty.Id.String())
		if err != nil {
			return
		}
		tokenEty = entities.Token{
			TokenableType: reflecth.GetTypeName(userEty),
			TokenableId:   userEty.Id,
			Type:          token_types.ResetPassword,
			Body:          strconv.Itoa(numberh.Random(100000, 999999)),
			UsedAt:        sql.NullTime{Time: time.Time{}, Valid: false},
			ExpiredAt:     sql.NullTime{Time: time.Now().Add(time.Hour / 2), Valid: true}, // give 30 mins expiration
		}
		_, err = service.tokenRepository.Create(ctx, &tokenEty)
		if err != nil {
			return
		}
	}

	// Send Token to User's email
	err = mailh.Send(os.Getenv("MAIL_SENDER"),
		"OTP | Reset Password",
		"Your reset password OTP is "+tokenEty.Body,
		data.Email,
	)
	if err != nil {
		return
	}

	return nil
}

// ResetPassword resets User's password if the given otp is valid
func (service *AuthService) ResetPassword(ctx context.Context, data authReqs.ResetPassword) (err error) {
	var (
		userEty  entities.User
		tokenEty entities.Token
		errChan  = chanh.Make[error](nil, 1)
	)

	userEty, err = service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		return
	}
	if userEty.Id == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}

	tokenEty, err = service.tokenRepository.FindByTokenableAndType(
		ctx,
		reflecth.GetTypeName(userEty),
		userEty.Id.String(),
		token_types.ResetPassword,
	)
	if err != nil {
		return
	}
	if tokenEty.Id == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}
	if tokenEty.UsedAt.Valid {
		err = exceptions.NewHTTPError(http.StatusUnprocessableEntity,
			exceptions.NewValidationError("otp", "OTP has been used / invalid OTP"))
		return
	}
	if tokenEty.IsExpired() {
		err = exceptions.NewHTTPError(http.StatusUnprocessableEntity,
			exceptions.NewValidationError("otp", "OTP is expired"))
		return
	}
	if tokenEty.Body != data.Otp {
		err = exceptions.NewHTTPError(http.StatusUnprocessableEntity,
			exceptions.NewValidationError("otp", "wrong or invalid OTP"))
		return
	}

	service.waitGroup.Add(2)

	go func() { // goroutine for update token
		defer service.waitGroup.Done()

		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		service.mutex.Lock()
		tokenEty.UsedAt = sql.NullTime{Time: time.Now(), Valid: true}
		_, errChanVal = service.tokenRepository.UpdateById(ctx, &tokenEty)
		service.mutex.Unlock()
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
	}()

	go func() { // goroutine for update user
		defer service.waitGroup.Done()

		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		passwordBytes, errChanVal := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

		service.mutex.Lock()
		userEty.Password = string(passwordBytes)
		_, errChanVal = service.userRepository.UpdateById(ctx, &userEty)
		service.mutex.Unlock()

		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
	}()

	service.waitGroup.Wait()

	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}

	return nil
}
