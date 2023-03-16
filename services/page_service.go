package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	media_coll_names "github.com/notefan-golang/enums/media/collection_names"
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
