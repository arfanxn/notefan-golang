package repositories

import (
	"context"
)

type RepositoryContract[T any] interface {
	Get(ctx context.Context) ([]T, error)
	Find(ctx context.Context, id string) (T, error)
	Create(ctx context.Context, model T) (T, error)
	Update(ctx context.Context, model T) (T, error)
	Delete(ctx context.Context, id string) error
}
