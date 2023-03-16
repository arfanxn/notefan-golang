package services

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	media_coll_names "github.com/notefan-golang/enums/media/collection_names"
	media_disks "github.com/notefan-golang/enums/media/disks"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/helpers/synch"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/page_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/page_ress"
	"github.com/notefan-golang/models/responses/pagination_ress"
	"github.com/notefan-golang/repositories"
)

type PageService struct {
	repository      *repositories.PageRepository
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.RWMutex
}

func NewPageService(
	repository *repositories.PageRepository,
	mediaRepository *repositories.MediaRepository,
) *PageService {
	return &PageService{
		repository:      repository,
		mediaRepository: mediaRepository,
		waitGroup:       new(sync.WaitGroup),
		mutex:           sync.RWMutex{},
	}
}

// GetBySpace get pages by space id and parse it to slice of Page Response
func (service *PageService) GetBySpace(ctx context.Context, data page_reqs.GetBySpace) (
	paginationRes pagination_ress.Pagination[page_ress.Page], err error) {
	service.repository.Query.Limit = data.PerPage
	service.repository.Query.Offset = (data.Page - 1) * int64(data.PerPage)
	service.repository.Query.Keyword = data.Keyword
	pageEtys, err := service.repository.GetBySpaceId(ctx, data.SpaceId)
	if err != nil {
		return
	}
	if len(pageEtys) == 0 { // immediately return if no Pages found
		return
	}

	// Retrieve Icon of each Pages
	pageIconMediaEtys, err := service.mediaRepository.GetByModelsAndCollectionNames(ctx,
		sliceh.Map(pageEtys, func(pageEty entities.Page) entities.Media {
			return entities.Media{
				ModelType:      reflecth.GetTypeName(pageEty),
				ModelId:        pageEty.Id,
				CollectionName: media_coll_names.Icon,
			}
		})...,
	)

	// Load Icon for each Pages
	for index, pageEty := range pageEtys {
		pageRes := page_ress.FillFromEntity(pageEty)
		paginationRes.Items = append(paginationRes.Items, pageRes)

		service.waitGroup.Add(1)
		go func(pageRes *page_ress.Page) {
			defer service.waitGroup.Done()
			iconMediaEtys := sliceh.Filter(pageIconMediaEtys, func(media entities.Media) bool {
				return pageRes.Id == media.ModelId.String()
			})
			if len(iconMediaEtys) == 0 {
				return
			}

			service.mutex.Lock()
			defer service.mutex.Unlock()
			pageRes.Icon = media_ress.FillFromEntity(iconMediaEtys[0])
		}(&paginationRes.Items[index])
	}

	service.waitGroup.Wait()

	paginationRes.Total = len(paginationRes.Items)
	return paginationRes, nil
}

// Find finds Page by the given data
func (service *PageService) Find(ctx context.Context, data page_reqs.Action) (
	pageRes page_ress.Page, err error) {
	var (
		pageEty entities.Page
		errChan = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan) // defer close channel

	service.waitGroup.Add(2)

	go func() { // groutine load Page
		defer service.waitGroup.Done()

		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}

		ety, errChanVal := service.repository.Find(ctx, data.PageId)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		if ety.Id == uuid.Nil {
			errChan <- exceptions.HTTPNotFound
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		iconMediaRes := pageRes.Icon
		pageRes = page_ress.FillFromEntity(ety)
		pageRes.Icon = iconMediaRes
	}()

	go func() { // groutine load Page's icon
		defer service.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		mediaEtys, errChanVal := service.mediaRepository.GetByModelsAndCollectionNames(ctx,
			entities.Media{
				ModelType:      reflecth.GetTypeName(pageEty),
				ModelId:        uuid.MustParse(data.PageId),
				CollectionName: media_coll_names.Icon,
			},
			entities.Media{
				ModelType:      reflecth.GetTypeName(pageEty),
				ModelId:        uuid.MustParse(data.PageId),
				CollectionName: media_coll_names.Cover,
			},
		)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		for _, mediaEty := range mediaEtys {
			service.mutex.Lock()
			switch mediaEty.CollectionName {
			case media_coll_names.Icon:
				pageRes.Icon = media_ress.FillFromEntity(mediaEty)
			case media_coll_names.Cover:
				pageRes.Cover = media_ress.FillFromEntity(mediaEty)
			}
			service.mutex.Unlock()
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

	return pageRes, nil
}

// Create creates Page from the given data
func (service *PageService) Create(ctx context.Context, data page_reqs.Create) (
	pageRes page_ress.Page, err error) {
	var (
		pageEty   entities.Page
		mediaEtys []*entities.Media
		errChan   = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)
	// Preapre Page entity for Page creation
	pageEty.SpaceId = uuid.MustParse(data.SpaceId)
	pageEty.Title = data.Title
	pageEty.Order = data.Order
	_, err = service.repository.Create(ctx, &pageEty) // create
	if err != nil {
		return
	}
	// Fill Page response/resource
	pageRes = page_ress.FillFromEntity(pageEty)
	service.waitGroup.Add(2)
	go func() {
		defer service.waitGroup.Done()
		if data.Icon == nil || !data.Icon.IsProvided() {
			return
		}
		var iconMediaEty entities.Media
		iconMediaEty.ModelType = reflecth.GetTypeName(pageEty)
		iconMediaEty.ModelId = pageEty.Id
		iconMediaEty.CollectionName = media_coll_names.Icon
		iconMediaEty.Disk = media_disks.Public
		iconMediaEty.FileName = filepath.Base(data.Icon.Name)
		iconMediaEty.File = data.Icon
		service.mutex.Lock()
		defer service.mutex.Unlock()
		mediaEtys = append(mediaEtys, &iconMediaEty)
	}()
	go func() {
		defer service.waitGroup.Done()
		if data.Cover == nil || !data.Cover.IsProvided() {
			return
		}
		var coverMediaEty entities.Media
		coverMediaEty.ModelType = reflecth.GetTypeName(pageEty)
		coverMediaEty.ModelId = pageEty.Id
		coverMediaEty.CollectionName = media_coll_names.Cover
		coverMediaEty.Disk = media_disks.Public
		coverMediaEty.FileName = filepath.Base(data.Cover.Name)
		coverMediaEty.File = data.Cover
		service.mutex.Lock()
		defer service.mutex.Unlock()
		mediaEtys = append(mediaEtys, &coverMediaEty)

	}()
	service.waitGroup.Wait()
	if len(mediaEtys) != 0 {
		_, err = service.mediaRepository.Insert(ctx, mediaEtys...)
		if err != nil {
			return
		}
	}
	for _, mediaEty := range mediaEtys {
		switch mediaEty.CollectionName {
		case media_coll_names.Icon:
			pageRes.Icon = media_ress.FillFromEntity(*mediaEty)
		case media_coll_names.Cover:
			pageRes.Cover = media_ress.FillFromEntity(*mediaEty)
		}
	}
	return
}

// Update updates Page by the given data
func (service *PageService) Update(ctx context.Context, data page_reqs.Update) (
	pageRes page_ress.Page, err error) {
	var (
		pageEty entities.Page
		errChan = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)
	pageEty, err = service.repository.Find(ctx, data.PageId)
	if err != nil {
		return
	}
	// Fill Page response/resource
	pageRes = page_ress.FillFromEntity(pageEty)
	service.waitGroup.Add(3)
	go func() { // goroutine for update Page
		defer service.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety := pageEty
		ety.Title = data.Title
		ety.Order = data.Order
		_, errChanVal = service.repository.UpdateById(ctx, &ety)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		pageEty = ety
		iconMediaRes, coverMediaRes := pageRes.Icon, pageRes.Cover
		pageRes = page_ress.FillFromEntity(ety)
		pageRes.Icon, pageRes.Cover = iconMediaRes, coverMediaRes
	}()
	go func() { // goroutine for update Page's Icon if exists
		defer service.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		// if icon is nil or not provided return immediately
		if data.Icon == nil || !data.Icon.IsProvided() {
			return
		}
		// Get Page's Icon that will be used for the update operation
		ety, errChanVal := service.mediaRepository.FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(entities.Page{}), data.PageId, media_coll_names.Icon,
		)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		if ety.Id == uuid.Nil { // create if space does not have Icon
			ety.ModelType = reflecth.GetTypeName(entities.Page{})
			ety.ModelId = uuid.MustParse(data.PageId)
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
			pageRes.Icon = media_ress.FillFromEntity(ety)
			return
		}
		// Update the Page's icon
		ety.FileName = data.Icon.Name
		ety.File = data.Icon
		_, errChanVal = service.mediaRepository.UpdateById(ctx, &ety)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		pageRes.Icon = media_ress.FillFromEntity(ety)
	}()
	go func() { // goroutine for update Page's Cover if exists
		defer service.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		// if Cover is nil or not provided return immediately
		if data.Cover == nil || !data.Cover.IsProvided() {
			return
		}
		// Get Page's Cover that will be used for the update operation
		ety, errChanVal := service.mediaRepository.FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(entities.Page{}), data.PageId, media_coll_names.Cover,
		)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		if ety.Id == uuid.Nil { // create if space does not have Cover
			ety.ModelType = reflecth.GetTypeName(entities.Page{})
			ety.ModelId = uuid.MustParse(data.PageId)
			ety.CollectionName = media_coll_names.Cover
			ety.Disk = media_disks.Public
			ety.FileName = data.Cover.Name
			ety.File = data.Cover
			_, errChanVal = service.mediaRepository.Create(ctx, &ety)
			if errChanVal != nil {
				errChan <- errChanVal
				return
			}
			service.mutex.Lock()
			defer service.mutex.Unlock()
			pageRes.Cover = media_ress.FillFromEntity(ety)
			return
		}
		// Update the Page's Cover
		ety.FileName = data.Cover.Name
		ety.File = data.Cover
		_, errChanVal = service.mediaRepository.UpdateById(ctx, &ety)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		pageRes.Cover = media_ress.FillFromEntity(ety)
	}()
	service.waitGroup.Wait() // waits
	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}
	return
}

// Delete deletes Page by the given data
func (service *PageService) Delete(ctx context.Context, data page_reqs.Action) error {
	errChan := synch.MakeChanWithValue[error](nil, 1)
	defer close(errChan)
	pageEty, err := service.repository.Find(ctx, data.PageId)
	if err != nil {
		return err
	}
	if pageEty.Id == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	service.waitGroup.Add(2)
	go func() { // goroutine for delete Space
		defer service.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		_, errChanVal = service.repository.DeleteByIds(ctx, pageEty.Id.String())
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
	}()
	go func() { // goroutine for delete Space's Medias
		defer service.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		mediaEtys, errChanVal := service.mediaRepository.GetByModelsAndCollectionNames(ctx,
			entities.Media{
				ModelType:      reflecth.GetTypeName(entities.Page{}),
				ModelId:        uuid.MustParse(data.PageId),
				CollectionName: media_coll_names.Icon,
			},
			entities.Media{
				ModelType:      reflecth.GetTypeName(entities.Page{}),
				ModelId:        uuid.MustParse(data.PageId),
				CollectionName: media_coll_names.Cover,
			},
		)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		var mediaIds []string
		for _, mediaEty := range mediaEtys {
			mediaIds = append(mediaIds, mediaEty.Id.String())
			service.waitGroup.Add(1)
			go func(mediaEty entities.Media) {
				defer service.waitGroup.Done()
				errChanVal = mediaEty.RemoveDirFile()
				if errChanVal != nil {
					errChan <- errChanVal
					return
				}
			}(mediaEty)
		}
		_, errChanVal = service.mediaRepository.DeleteByIds(ctx, mediaIds...)
		if errChanVal != nil {
			errChan <- errChanVal
			return
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
