package ctrl

import (
	"context"
	"time"
)

type appRepo interface {
	itemRepo
	categoryRepo
	promotionRepo
	favoriteRepo
}

type CacheRepo interface {
	GetCode(ctx context.Context, key string) (int, error)
	GetToStruct(ctx context.Context, key string, dest any) error

	Set(ctx context.Context, t time.Duration, key string, val any) error
	Delete(ctx context.Context, key string) error
	Close()

	InvalidateKeysByPattern(ctx context.Context, pattern string) error
}

type Controller struct {
	repo  appRepo
	cache CacheRepo
}

func New(repo appRepo, cache CacheRepo) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
	}
}
