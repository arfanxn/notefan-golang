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
	"github.com/notefan-golang/helpers/errorh"
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

func (service *AuthService) Login(ctx context.Context, data authReqs.Login) (authRess.Login, error) {
	user, err := service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil { // err not nil == user not found, return exception HTTPAuthLoginFailed
		errorh.Log(err)
		return authRess.Login{}, exceptions.HTTPAuthLoginFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil { // err not nil == password doesnt match, return exception HTTPAuthLoginFailed
		errorh.Log(err)
		return authRess.Login{}, exceptions.HTTPAuthLoginFailed
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
	errorh.LogPanic(err) // panic if token generation failed

	return authRess.Login{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}

// Register registers the given user
func (service *AuthService) Register(ctx context.Context, data authReqs.Register) (user_ress.User, error) {
	// Hash the user password
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	errorh.LogPanic(err) // panic if password hashing failed

	var (
		userEty  entities.User
		mediaEty entities.Media
	)

	service.waitGroup.Add(2)

	go func() { // goroutine for creating user
		defer service.waitGroup.Done()

		service.mutex.Lock()
		// parse request to entity
		userEty = entities.User{
			Name:     data.Name,
			Email:    data.Email,
			Password: string(password),
		}
		service.mutex.Unlock()

		// Save user into Database
		_, err = service.userRepository.Create(ctx, &userEty)
		errorh.LogPanic(err) // panic if save into db failed
	}()

	go func() { // goroutine for creating user's avatar
		defer service.waitGroup.Done()

		// the user default avatar file path
		defaultAvatarFilePath := "./public/placeholders/images/default-avatar.jpg"
		// open default user avatar file
		defaultAvatarBytes, err := os.ReadFile(defaultAvatarFilePath)
		errorh.LogPanic(err) // panic if failed to read file

		service.mutex.Lock()
		// Prepare media entity for user avatar
		mediaEty = entities.Media{
			ModelType:      reflecth.GetTypeName(userEty),
			ModelId:        userEty.Id,
			CollectionName: media_collnames.Avatar,
			Disk:           media_disks.Public,
			FileName:       filepath.Base(defaultAvatarFilePath),
			File:           file_reqs.NewFromBytes(defaultAvatarBytes),
		}
		service.mutex.Unlock()

		// Save media into Database
		_, err = service.mediaRepository.Create(ctx, &mediaEty)
		errorh.LogPanic(err) // panic if save into db failed
	}()

	service.waitGroup.Wait()

	userRes := user_ress.FillFromEntity(userEty)
	userRes.Avatar = media_ress.FillFromEntity(mediaEty)

	// Return the created user and nil
	return userRes, nil
}
