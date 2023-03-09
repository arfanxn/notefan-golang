package services

import (
	"context"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/requests/space_reqs"
	"github.com/notefan-golang/models/responses/pagination_ress"
	"github.com/notefan-golang/models/responses/space_ress"
	"github.com/notefan-golang/repositories"
)

type SpaceService struct {
	repository *repositories.SpaceRepository
}

func NewSpaceService(
	repository *repositories.SpaceRepository,
) *SpaceService {
	return &SpaceService{repository: repository}
}

// GetByUserId get spaces by user id and parse it to slice of Space Response
func (service *SpaceService) GetByUserId(ctx context.Context, data space_reqs.GetByUser) (pagination_ress.Pagination[space_ress.Space], error) {
	service.repository.Query.Limit = data.PerPage
	service.repository.Query.Offset = (data.Page - 1) * int64(data.PerPage)
	spaceEtys, err := service.repository.GetByUserId(ctx, data.UserId)
	errorh.Panic(err)
	spaceRess := space_ress.FillFromEntities(spaceEtys)
	pagination := pagination_ress.Make[space_ress.Space]()
	pagination.SetItems(spaceRess)
	return pagination, nil
}

// Find finds space by the given request id
func (service *SpaceService) Find(ctx context.Context, data common_reqs.UUID) (space_ress.Space, error) {
	spaceEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err)
	return space_ress.FillFromEntity(spaceEty), nil
}

// Update updates space by the given request id
func (service *SpaceService) Update(ctx context.Context, data space_reqs.Space) (space_ress.Space, error) {
	spaceEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err) // panic if not found

	if data.Name != "" {
		spaceEty.Name = data.Name
	}
	if data.Description != "" {
		spaceEty.Description = data.Description
	}
	if data.Domain != "" {
		spaceEty.Domain = data.Domain
	}

	_, err = service.repository.UpdateById(ctx, &spaceEty)
	errorh.LogPanic(err)

	return space_ress.FillFromEntity(spaceEty), nil
}

// Delete deletes space by the given request id
func (service *SpaceService) Delete(ctx context.Context, data common_reqs.UUID) error {
	spaceEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err) // panic if not found

	_, err = service.repository.DeleteByIds(ctx, spaceEty.Id.String())
	return err
}
