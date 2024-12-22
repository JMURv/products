package ctrl

import (
	"context"
	"time"
)

type AppRepo interface {
	itemRepo
	categoryRepo
	promotionRepo
	favoriteRepo
}

type Discovery interface {
	Register() error
	Deregister() error
	FindServiceByName(ctx context.Context, name string) (string, error)
}

type CacheService interface {
	GetCode(ctx context.Context, key string) (int, error)
	GetToStruct(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, t time.Duration, key string, val any) error
	Delete(ctx context.Context, key string) error
	Close()
	InvalidateKeysByPattern(ctx context.Context, pattern string) error
}

type Controller struct {
	repo  AppRepo
	cache CacheService
}

func New(repo AppRepo, cache CacheService) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
	}
}
