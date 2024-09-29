package ctrl

import (
	"context"
	"github.com/JMURv/par-pro/products/internal/ctrl/etc"
	"github.com/JMURv/par-pro/products/internal/ctrl/seo"
	"time"
)

type appRepo interface {
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
	seo   seo.SEOCtrl
	etc   etc.EtcCtrl
}

func New(repo appRepo, cache CacheRepo, seo seo.SEOCtrl, etc etc.EtcCtrl) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
		seo:   seo,
		etc:   etc,
	}
}
