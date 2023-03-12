package services

import (
	"context"
	"errors"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	coll_names "github.com/notefan-golang/enums/media/collection_names"
	media_disks "github.com/notefan-golang/enums/media/disks"
	roleNames "github.com/notefan-golang/enums/role/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/requests/space_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/pagination_ress"
	"github.com/notefan-golang/models/responses/space_ress"
	"github.com/notefan-golang/repositories"
	"gopkg.in/guregu/null.v4"
)

type SpaceService struct {
	repository      *repositories.SpaceRepository
	ursRepository   *repositories.UserRoleSpaceRepository
	roleRepository  *repositories.RoleRepository
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.RWMutex
}

func NewSpaceService(
	repository *repositories.SpaceRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
	roleRepository *repositories.RoleRepository,
	mediaRepository *repositories.MediaRepository,
) *SpaceService {
	return &SpaceService{
		repository:      repository,
		ursRepository:   ursRepository,
		roleRepository:  roleRepository,
		mediaRepository: mediaRepository,
		waitGroup:       new(sync.WaitGroup),
		mutex:           sync.RWMutex{},
	}
}

// GetByUser get spaces by user id and parse it to slice of Space Response
func (service *SpaceService) GetByUser(ctx context.Context, data space_reqs.GetByUser) (pagination_ress.Pagination[space_ress.Space], error) {
	service.repository.Query.Limit = data.PerPage
	service.repository.Query.Offset = (data.Page - 1) * int64(data.PerPage)
	service.repository.Query.Keyword = data.Keyword
	spaceEtys, err := service.repository.GetByUserId(ctx, data.UserId)
	errorh.Panic(err)
	var spaceRess []space_ress.Space
	for _, spaceEty := range spaceEtys {
		var spaceRes space_ress.Space
		spaceRes.Id = spaceEty.Id.String()
		spaceRes.Name = spaceEty.Name
		spaceRes.Description = spaceEty.Description
		spaceRes.Domain = spaceEty.Domain
		spaceRes.CreatedAt = spaceEty.CreatedAt
		spaceRes.UpdatedAt = null.NewTime(spaceEty.UpdatedAt.Time, spaceEty.UpdatedAt.Valid)
		spaceRess = append(spaceRess, spaceRes)
	}

	spaceIconMediaEtys, err := service.mediaRepository.
		GetByModelsAndCollectionNames(ctx,
			sliceh.Map(spaceEtys, func(spaceEty entities.Space) entities.Media {
				return entities.Media{
					ModelType:      reflecth.GetTypeName(spaceEty),
					ModelId:        spaceEty.Id,
					CollectionName: coll_names.Icon,
				}
			})...,
		)

	// Load Icon for each Spaces
	for i := 0; i < len(spaceRess); i++ {
		spaceRes := &spaceRess[i]

		service.waitGroup.Add(1)
		go func(spaceRes *space_ress.Space) {
			defer service.waitGroup.Done()
			iconMediaEtys := sliceh.Filter(spaceIconMediaEtys, func(media entities.Media) bool {
				return spaceRes.Id == media.ModelId.String()
			})
			if len(iconMediaEtys) == 0 {
				return
			}

			service.mutex.Lock()
			defer service.mutex.Unlock()
			spaceRes.Icon = media_ress.FillFromEntity(iconMediaEtys[0])
		}(spaceRes)
	}

	service.waitGroup.Wait()

	pagination := pagination_ress.Make[space_ress.Space]()
	pagination.SetItems(spaceRess)
	return pagination, nil
}

// Find finds space by the given request id
func (service *SpaceService) Find(ctx context.Context, data common_reqs.UUID) (spaceRes space_ress.Space, err error) {
	var spaceEty entities.Space

	service.waitGroup.Add(2)

	go func() { // groutine for find space by id
		defer service.waitGroup.Done()

		service.mutex.Lock()
		defer service.mutex.Unlock()
		spaceEty, err = service.repository.Find(ctx, data.Id)
		errorh.Panic(err)
		iconMediaRes := spaceRes.Icon
		spaceRes = space_ress.FillFromEntity(spaceEty)
		spaceRes.Icon = iconMediaRes
	}()

	go func() {
		defer service.waitGroup.Done()

		iconMediaEty, err := service.mediaRepository.FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(spaceEty), data.Id, coll_names.Icon,
		)
		if errors.Is(err, exceptions.HTTPNotFound) {
			return // return from goroutine if space's icon media is not found
		}
		errorh.LogPanic(err)

		service.mutex.Lock()
		defer service.mutex.Unlock()
		spaceRes.Icon = media_ress.FillFromEntity(iconMediaEty)
	}()

	service.waitGroup.Wait()

	return spaceRes, nil
}

// Create creates space from the given data
func (service *SpaceService) Create(ctx context.Context, data space_reqs.Create) (
	spaceRes space_ress.Space, err error) {
	var (
		spaceEty     entities.Space
		iconMediaEty entities.Media
	)
	spaceEty.Name = data.Name
	spaceEty.Description = data.Description
	spaceEty.Domain = data.Domain

	_, err = service.repository.Create(ctx, &spaceEty)
	errorh.LogPanic(err)

	// Fill Space response/resource
	spaceRes = space_ress.FillFromEntity(spaceEty)

	service.waitGroup.Add(2)

	go func() { // goroutine for save space's icon if exists
		defer service.waitGroup.Done()

		// if icon is nil or not provided return immediately
		if data.Icon == nil || !data.Icon.IsProvided() {
			return
		}

		// Prepare Media entity
		var mediaEty entities.Media
		mediaEty.ModelType = reflecth.GetTypeName(spaceEty)
		mediaEty.ModelId = spaceEty.Id
		mediaEty.CollectionName = coll_names.Icon
		mediaEty.Disk = media_disks.Public
		mediaEty.FileName = filepath.Base(data.Icon.Name)
		mediaEty.File = data.Icon

		// Create the media entity
		_, err = service.mediaRepository.Create(ctx, &mediaEty)
		errorh.LogPanic(err)

		service.mutex.Lock()
		defer service.mutex.Unlock()
		iconMediaEty = mediaEty
		spaceRes.Icon = media_ress.FillFromEntity(iconMediaEty)
	}()

	go func() { // goroutine for create relationship between space, user and role
		defer service.waitGroup.Done()

		// Get Space ownership role
		roleEty, err := service.roleRepository.FindByName(ctx, roleNames.SpaceOwner)
		errorh.Panic(err)

		ursEty := entities.UserRoleSpace{
			UserId:  uuid.MustParse(data.UserId),
			RoleId:  roleEty.Id,
			SpaceId: spaceEty.Id,
		}
		_, err = service.ursRepository.Create(ctx, &ursEty)
		errorh.Panic(err)
	}()

	service.waitGroup.Wait()

	return
}

// Update updates space by the given request id
func (service *SpaceService) Update(ctx context.Context, data space_reqs.Update) (
	spaceRes space_ress.Space, err error) {
	spaceEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err) // panic if not found

	service.waitGroup.Add(2)

	go func() { // go routine for update Space entity
		defer service.waitGroup.Done()

		if data.Name != "" {
			spaceEty.Name = data.Name
		}
		if data.Description != "" {
			spaceEty.Description = data.Description
		}
		if data.Domain != "" {
			spaceEty.Domain = data.Domain
		}

		_, err := service.repository.UpdateById(ctx, &spaceEty)
		errorh.LogPanic(err)

		service.mutex.Lock()
		defer service.mutex.Unlock()
		iconMediaRes := spaceRes.Icon // get the icon incase of the icon is already set in the other goroutine
		spaceRes = space_ress.FillFromEntity(spaceEty)
		spaceRes.Icon = iconMediaRes
	}()

	go func() { // go routine for update Space's icon if exists
		defer service.waitGroup.Done()

		// if icon is nil or not provided return immediately
		if data.Icon == nil || !data.Icon.IsProvided() {
			return
		}

		// Get Space's Icon that will be used for the update operation
		mediaEty, err := service.mediaRepository.FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(spaceEty), spaceEty.Id.String(), coll_names.Icon,
		)
		errorh.LogPanic(err)

		mediaEty.FileName = data.Icon.Name
		mediaEty.File = data.Icon

		// Update the Space's icon
		_, err = service.mediaRepository.UpdateById(ctx, &mediaEty)
		errorh.LogPanic(err)

		service.mutex.Lock()
		defer service.mutex.Unlock()
		spaceRes.Icon = media_ress.FillFromEntity(mediaEty)
	}()

	service.waitGroup.Wait()

	return spaceRes, nil
}

// Delete deletes space by the given request id
func (service *SpaceService) Delete(ctx context.Context, data common_reqs.UUID) error {
	spaceEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err) // panic if not found

	_, err = service.repository.DeleteByIds(ctx, spaceEty.Id.String())
	return err
}
