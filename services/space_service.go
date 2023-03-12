package services

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	media_coll_names "github.com/notefan-golang/enums/media/collection_names"
	media_disks "github.com/notefan-golang/enums/media/disks"
	roleNames "github.com/notefan-golang/enums/role/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/helpers/synch"
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
func (service *SpaceService) GetByUser(ctx context.Context, data space_reqs.GetByUser) (
	paginationRes pagination_ress.Pagination[space_ress.Space], err error) {
	service.repository.Query.Limit = data.PerPage
	service.repository.Query.Offset = (data.Page - 1) * int64(data.PerPage)
	service.repository.Query.Keyword = data.Keyword
	spaceEtys, err := service.repository.GetByUserId(ctx, data.UserId)
	if err != nil {
		return
	}
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

	if len(spaceEtys) == 0 { // immediately return if no Spaces found
		return
	}

	spaceIconMediaEtys, err := service.mediaRepository.
		GetByModelsAndCollectionNames(ctx,
			sliceh.Map(spaceEtys, func(spaceEty entities.Space) entities.Media {
				return entities.Media{
					ModelType:      reflecth.GetTypeName(spaceEty),
					ModelId:        spaceEty.Id,
					CollectionName: media_coll_names.Icon,
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

	paginationRes.SetItems(spaceRess)
	return paginationRes, nil
}

// Find finds space by the given request id
func (service *SpaceService) Find(ctx context.Context, data common_reqs.UUID) (
	spaceRes space_ress.Space, err error) {
	var (
		spaceEty entities.Space
		errChan  = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan) // defer close channel

	service.waitGroup.Add(2)

	go func() { // groutine for find space by id
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		ety, errChanVal := service.repository.Find(ctx, data.Id)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		if ety.Id == uuid.Nil {
			// return exceptions.HTTPNotFound if space not found
			errChan <- exceptions.HTTPNotFound
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		iconMediaRes := spaceRes.Icon
		spaceRes = space_ress.FillFromEntity(ety)
		spaceRes.Icon = iconMediaRes
	}()

	go func() { // groutine for get space's icon
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		iconMediaEty, errChanVal := service.mediaRepository.FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(spaceEty), data.Id, media_coll_names.Icon,
		)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}

		if iconMediaEty.Id == uuid.Nil {
			return // return from goroutine if Space does not have Icon
		}

		service.mutex.Lock()
		defer service.mutex.Unlock()
		spaceRes.Icon = media_ress.FillFromEntity(iconMediaEty)
	}()

	service.waitGroup.Wait()

	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}

	return spaceRes, nil
}

// Create creates space from the given data
func (service *SpaceService) Create(ctx context.Context, data space_reqs.Create) (
	spaceRes space_ress.Space, err error) {
	var (
		spaceEty     entities.Space
		iconMediaEty entities.Media
		errChan      = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)
	spaceEty.Name = data.Name
	spaceEty.Description = data.Description
	spaceEty.Domain = data.Domain

	_, err = service.repository.Create(ctx, &spaceEty)
	if err != nil {
		return
	}

	// Fill Space response/resource
	spaceRes = space_ress.FillFromEntity(spaceEty)

	service.waitGroup.Add(2)

	go func() { // goroutine for save space's icon if exists
		defer service.waitGroup.Done()

		// if icon is nil or not provided return immediately
		if data.Icon == nil || !data.Icon.IsProvided() {
			return
		}

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		// Prepare Media entity
		var mediaEty entities.Media
		mediaEty.ModelType = reflecth.GetTypeName(spaceEty)
		mediaEty.ModelId = spaceEty.Id
		mediaEty.CollectionName = media_coll_names.Icon
		mediaEty.Disk = media_disks.Public
		mediaEty.FileName = filepath.Base(data.Icon.Name)
		mediaEty.File = data.Icon

		// Create the media entity
		_, errChanVal = service.mediaRepository.Create(ctx, &mediaEty)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}

		service.mutex.Lock()
		defer service.mutex.Unlock()
		iconMediaEty = mediaEty
		spaceRes.Icon = media_ress.FillFromEntity(iconMediaEty)
	}()

	go func() { // goroutine for create relationship between space, user and role
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		// Get Space ownership role
		roleEty, errChanVal := service.roleRepository.FindByName(ctx, roleNames.SpaceOwner)
		if errChanVal != nil {
			errChan <- errChanVal
		}

		ursEty := entities.UserRoleSpace{
			UserId:  uuid.MustParse(data.UserId),
			RoleId:  roleEty.Id,
			SpaceId: spaceEty.Id,
		}
		_, errChanVal = service.ursRepository.Create(ctx, &ursEty)
		if errChanVal != nil {
			errChan <- errChanVal
		}
	}()

	service.waitGroup.Wait()

	return
}

// Update updates space by the given request id
func (service *SpaceService) Update(ctx context.Context, data space_reqs.Update) (
	spaceRes space_ress.Space, err error) {
	var (
		spaceEty entities.Space
		errChan  = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)

	spaceEty, err = service.repository.Find(ctx, data.Id)
	if err != nil {
		return
	}
	if spaceEty.Id == uuid.Nil {
		// return exceptions.HTTPNotFound if space not found
		err = exceptions.HTTPNotFound
		return
	}

	service.waitGroup.Add(2)

	go func() { // goroutine for update Space entity
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		if data.Name != "" {
			spaceEty.Name = data.Name
		}
		if data.Description != "" {
			spaceEty.Description = data.Description
		}
		if data.Domain != "" {
			spaceEty.Domain = data.Domain
		}

		_, errChanVal = service.repository.UpdateById(ctx, &spaceEty)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}

		service.mutex.Lock()
		defer service.mutex.Unlock()
		iconMediaRes := spaceRes.Icon // get the icon incase of the icon is already set in the other goroutine
		spaceRes = space_ress.FillFromEntity(spaceEty)
		spaceRes.Icon = iconMediaRes
	}()

	go func() { // goroutine for update Space's icon if exists
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		// if icon is nil or not provided return immediately
		if data.Icon == nil || !data.Icon.IsProvided() {
			return
		}

		// Get Space's Icon that will be used for the update operation
		ety, errChanVal := service.mediaRepository.FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(spaceEty), spaceEty.Id.String(), media_coll_names.Icon,
		)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		if ety.Id == uuid.Nil { // create if space does not have Icon
			ety.ModelType = reflecth.GetTypeName(spaceEty)
			ety.ModelId = spaceEty.Id
			ety.CollectionName = media_coll_names.Icon
			ety.Disk = media_disks.Public
			ety.FileName = data.Icon.Name
			ety.File = data.Icon
			_, errChanVal = service.mediaRepository.Create(ctx, &ety)
			if errChanVal != nil {
				errChan <- errChanVal
				return
			}

			service.mutex.Lock()
			defer service.mutex.Unlock()
			spaceRes.Icon = media_ress.FillFromEntity(ety)
			return
		}

		// Update the Space's icon
		ety.FileName = data.Icon.Name
		ety.File = data.Icon
		_, errChanVal = service.mediaRepository.UpdateById(ctx, &ety)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}

		service.mutex.Lock()
		defer service.mutex.Unlock()
		spaceRes.Icon = media_ress.FillFromEntity(ety)
	}()

	service.waitGroup.Wait()

	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}

	return spaceRes, nil
}

// Delete deletes space by the given request id
func (service *SpaceService) Delete(ctx context.Context, data common_reqs.UUID) error {
	errChan := synch.MakeChanWithValue[error](nil, 1)
	spaceEty, err := service.repository.Find(ctx, data.Id)
	if err != nil {
		return err
	}
	if spaceEty.Id == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	service.waitGroup.Add(2)

	go func() { // goroutine for delete Space
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		_, errChanVal = service.repository.DeleteByIds(ctx, spaceEty.Id.String())
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
	}()

	go func() { // goroutine for delete Space's Icon
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		iconMediaEty, errChanVal := service.mediaRepository.FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(spaceEty), spaceEty.Id.String(), media_coll_names.Icon)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}

		// if iconMediaEty.Id is not equals to uuid.Nil (space has icon)
		// Delete Icon's file and delete the Icon entity
		if iconMediaEty.Id != uuid.Nil {
			errChanVal = iconMediaEty.RemoveDirFile()
			if errChanVal != nil {
				errChan <- errChanVal
				return
			}

			_, errChanVal := service.mediaRepository.DeleteByIds(ctx, iconMediaEty.Id.String())
			if errChanVal != nil {
				errChan <- errChanVal
				return
			}
		}
	}()

	service.waitGroup.Wait()

	if err != nil {
		return err
	}
	err = <-errChan
	if err != nil {
		return err
	}
	return nil
}
