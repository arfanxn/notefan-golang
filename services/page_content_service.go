package services

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	media_coll_names "github.com/notefan-golang/enums/media/collection_names"
	media_disks "github.com/notefan-golang/enums/media/disks"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/chanh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/file_reqs"
	pc_reqs "github.com/notefan-golang/models/requests/page_content_reqs"
	"github.com/notefan-golang/models/requests/query_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	pc_ress "github.com/notefan-golang/models/responses/page_content_ress"
	"github.com/notefan-golang/models/responses/pagination_ress"
	"github.com/notefan-golang/repositories"
)

type PageContentService struct {
	repository      *repositories.PageContentRepository
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.RWMutex
}

func NewPageContentService(
	repository *repositories.PageContentRepository,
	mediaRepository *repositories.MediaRepository,
) *PageContentService {
	return &PageContentService{
		repository:      repository,
		mediaRepository: mediaRepository,
		waitGroup:       new(sync.WaitGroup),
		mutex:           sync.RWMutex{},
	}
}

// GetByPage get page contents by page id and parse it to slice of PageContent Response
func (service *PageContentService) GetBypage(ctx context.Context, data pc_reqs.GetByPage) (
	paginationRes pagination_ress.Pagination[pc_ress.PageContent], err error) {
	service.repository.Query.Limit = data.PerPage
	service.repository.Query.Offset = (data.Page - 1) * int64(data.PerPage)
	service.repository.Query.Keyword = data.Keyword
	for _, orderBy := range data.OrderBys {
		keyAndVal := strings.Split(orderBy, "=")
		service.repository.Query.AddOrderBy(keyAndVal[0], keyAndVal[1])
	}
	pcEtys, err := service.repository.GetByPageId(ctx, data.PageId)
	service.repository.Query = query_reqs.Default() // reset query to default after retrieving
	if err != nil {
		return
	}
	if len(pcEtys) == 0 { // immediately return if no PageContents found
		return
	}
	// Retrieve Medias of each PageContents
	pcMediaEtys, err := service.mediaRepository.GetByModelsAndCollectionNames(ctx,
		sliceh.Map(pcEtys, func(pcEty entities.PageContent) entities.Media {
			return entities.Media{
				ModelType:      reflecth.GetTypeName(pcEty),
				ModelId:        pcEty.Id,
				CollectionName: media_coll_names.ContentMedia,
			}
		})...,
	)
	// Load Medias for each PageContents
	for index, pcEty := range pcEtys {
		pcRes := pc_ress.FillFromEntity(pcEty)
		paginationRes.Items = append(paginationRes.Items, pcRes)
		service.waitGroup.Add(1)
		go func(pcRes *pc_ress.PageContent) {
			defer service.waitGroup.Done()
			mediaEtys := sliceh.Filter(pcMediaEtys, func(media entities.Media) bool {
				return pcRes.Id == media.ModelId.String()
			})
			if len(mediaEtys) == 0 {
				return
			}

			service.mutex.Lock()
			defer service.mutex.Unlock()
			pcRes.Medias = sliceh.Map(mediaEtys, func(mediaEty entities.Media) media_ress.Media {
				return media_ress.FillFromEntity(mediaEty)
			})
		}(&paginationRes.Items[index])
	}
	service.waitGroup.Wait()
	paginationRes.Total = len(paginationRes.Items)
	return paginationRes, nil
}

// Find finds PageContent by the given data
func (service *PageContentService) Find(ctx context.Context, data pc_reqs.Action) (
	pcRes pc_ress.PageContent, err error) {
	var (
		errChan = chanh.Make[error](nil, 1)
	)
	defer close(errChan) // defer close channel
	service.waitGroup.Add(2)
	go func() { // goroutine for load PageContent
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := service.repository.Find(ctx, data.PageContentId)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if ety.Id == uuid.Nil {
			errChan <- exceptions.HTTPNotFound
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		pcMediaRess := pcRes.Medias
		pcRes = pc_ress.FillFromEntity(ety)
		pcRes.Medias = pcMediaRess
	}()

	go func() { // goroutine for load PageContent's Medias
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		etys, errChanVal := service.mediaRepository.GetByModelsAndCollectionNames(ctx,
			entities.Media{
				ModelType:      reflecth.GetTypeName(entities.PageContent{}),
				ModelId:        uuid.MustParse(data.PageContentId),
				CollectionName: media_coll_names.ContentMedia,
			},
		)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		mediaRess := sliceh.Map(etys, func(ety entities.Media) media_ress.Media {
			return media_ress.FillFromEntity(ety)
		})
		pcRes.Medias = append(pcRes.Medias, mediaRess...)
	}()
	service.waitGroup.Wait()
	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}
	err = nil
	return
}

// Create creates PageContent from the given data
func (service *PageContentService) Create(ctx context.Context, data pc_reqs.Create) (
	pcRes pc_ress.PageContent, err error) {
	var (
		pcEty     entities.PageContent
		mediaEtys []*entities.Media
	)
	// Preapre PageContent entity for PageContent creation
	pcEty.PageId = uuid.MustParse(data.PageId)
	pcEty.Type = data.Type
	pcEty.Order = data.Order
	pcEty.Body = data.Body
	_, err = service.repository.Create(ctx, &pcEty) // create
	if err != nil {
		return
	}
	// Fill PageContent response/resource
	pcRes = pc_ress.FillFromEntity(pcEty)
	// Parse PageContent's Medias Request to Entity if exists
	for _, mediaFileReq := range data.Medias {
		service.waitGroup.Add(1)
		go func(fileReq *file_reqs.File) {
			defer service.waitGroup.Done()
			var ety entities.Media
			ety.ModelType = reflecth.GetTypeName(pcEty)
			ety.ModelId = pcEty.Id
			ety.CollectionName = media_coll_names.ContentMedia
			ety.Disk = media_disks.Public
			ety.FileName = filepath.Base(fileReq.Name)
			ety.File = fileReq
			service.mutex.Lock()
			defer service.mutex.Unlock()
			mediaEtys = append(mediaEtys, &ety)
		}(mediaFileReq)
	}
	service.waitGroup.Wait()
	if len(mediaEtys) != 0 {
		_, err = service.mediaRepository.Insert(ctx, mediaEtys...)
		if err != nil {
			return
		}
		mediaRess := sliceh.Map(mediaEtys, func(ety *entities.Media) media_ress.Media {
			return media_ress.FillFromEntity(*ety)
		})
		pcRes.Medias = append(pcRes.Medias, mediaRess...)
	}
	return
}

// Update updates PageContent by the given data
func (service *PageContentService) Update(ctx context.Context, data pc_reqs.Update) (
	pcRes pc_ress.PageContent, err error) {
	var (
		mediaEtys []*entities.Media
		errChan   = chanh.Make[error](nil, 1)
	)
	defer close(errChan)
	service.waitGroup.Add(2)
	go func() {
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := service.repository.Find(ctx, data.PageContentId)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if ety.Id == uuid.Nil {
			errChanVal = exceptions.HTTPNotFound
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		// Update the PageContent
		ety.Type = data.Type
		ety.Order = data.Order
		ety.Body = data.Body
		_, errChanVal = service.repository.UpdateById(ctx, &ety)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		pcMediaRess := pcRes.Medias
		pcRes = pc_ress.FillFromEntity(ety)
		pcRes.Medias = pcMediaRess
	}()
	go func() { // goroutine for deleting the old PageContent's Medias
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		// if no Medias return immediately
		if len(data.Medias) == 0 {
			return
		}
		// Get PageContent's Medias that will be deleted
		etys, errChanVal := service.mediaRepository.GetByModelsAndCollectionNames(ctx, entities.Media{
			ModelType:      reflecth.GetTypeName(entities.PageContent{}),
			ModelId:        uuid.MustParse(data.PageContentId),
			CollectionName: media_coll_names.ContentMedia,
		})
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if len(etys) == 0 { // if PageContent doesn't have any Medias return immediately
			return
		}
		ids := sliceh.Map(etys, func(ety entities.Media) string {
			return ety.Id.String()
		})
		_, errChanVal = service.mediaRepository.DeleteByIds(ctx, ids...)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
	}()
	for _, mediaFileReq := range data.Medias {
		service.waitGroup.Add(1)
		go func(fileReq *file_reqs.File) {
			defer service.waitGroup.Done()
			var ety entities.Media
			ety.ModelType = reflecth.GetTypeName(entities.PageContent{})
			ety.ModelId = uuid.MustParse(data.PageContentId)
			ety.CollectionName = media_coll_names.ContentMedia
			ety.Disk = media_disks.Public
			ety.FileName = fileReq.Name
			ety.File = fileReq
			service.mutex.Lock()
			defer service.mutex.Unlock()
			mediaEtys = append(mediaEtys, &ety)
		}(mediaFileReq)
	}
	service.waitGroup.Wait()
	if err = <-errChan; err != nil {
		return
	}
	// insert/associate Medias with PageContent if exists
	if len(mediaEtys) != 0 {
		_, err = service.mediaRepository.Insert(ctx, mediaEtys...)
		if err != nil {
			return
		}
		pcRes.Medias = sliceh.Map(mediaEtys, func(ety *entities.Media) media_ress.Media {
			return media_ress.FillFromEntity(*ety)
		})
	}
	return
}

// Delete deletes PageContent by the given data
func (service *PageContentService) Delete(ctx context.Context, data pc_reqs.Action) error {
	errChan := chanh.Make[error](nil, 1)
	defer close(errChan)
	pcEty, err := service.repository.Find(ctx, data.PageContentId)
	if err != nil {
		return err
	}
	if pcEty.Id == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	service.waitGroup.Add(2)
	go func() { // goroutine for delete PageContent
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		_, errChanVal = service.repository.DeleteByIds(ctx, pcEty.Id.String())
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
	}()
	go func() { // goroutine for delete PageContent's Medias
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		mediaEtys, errChanVal := service.mediaRepository.GetByModelsAndCollectionNames(ctx,
			entities.Media{
				ModelType:      reflecth.GetTypeName(entities.PageContent{}),
				ModelId:        uuid.MustParse(data.PageContentId),
				CollectionName: media_coll_names.ContentMedia,
			},
		)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
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
					chanh.ReplaceVal(errChan, errChanVal)
					return
				}
			}(mediaEty)
		}
		_, errChanVal = service.mediaRepository.DeleteByIds(ctx, mediaIds...)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
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
