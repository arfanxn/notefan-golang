package services

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	media_collnames "github.com/notefan-golang/enums/media/collection_names"
	media_disks "github.com/notefan-golang/enums/media/disks"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/jwth"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
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
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.Mutex
}

func NewAuthService(
	userRepository *repositories.UserRepository,
	mediaRepository *repositories.MediaRepository) *AuthService {
	return &AuthService{
		userRepository:  userRepository,
		mediaRepository: mediaRepository,
		waitGroup:       new(sync.WaitGroup),
		mutex:           sync.Mutex{},
	}
}

func (service *AuthService) Login(ctx context.Context, data authReqs.Login) (loginRes authRess.Login, err error) {
	user, err := service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil { // err not nil == user not found, return exception HTTPAuthLoginFailed
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
	)

	service.waitGroup.Add(2)

	go func() { // goroutine for creating user
		defer service.waitGroup.Done()

		if err != nil {
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
		_, err = service.userRepository.Create(ctx, &userEty)
	}()

	go func() { // goroutine for creating user's avatar
		defer service.waitGroup.Done()

		if err != nil {
			return
		}

		// the user default avatar file path
		defaultAvatarFilePath := "./public/placeholders/images/default-avatar.jpg"

		service.mutex.Lock()
		defer service.mutex.Unlock()

		// open default user avatar file
		defaultAvatarBytes, err := os.ReadFile(defaultAvatarFilePath)
		if err != nil {
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
		_, err = service.mediaRepository.Create(ctx, &mediaEty)
	}()

	service.waitGroup.Wait()

	if err != nil {
		return
	}

	userRes = user_ress.FillFromEntity(userEty)
	userRes.Avatar = media_ress.FillFromEntity(mediaEty)

	// Return the created user and nil
	return userRes, nil
}
