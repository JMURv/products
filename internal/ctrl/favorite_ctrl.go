package ctrl

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const favoriteCacheKey = "favorite:%v"

type favoriteRepo interface {
	GetFavorites(ctx context.Context, uid uuid.UUID) ([]*model.Favorite, error)
	AddToFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) (*model.Favorite, error)
	RemoveFromFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) error
}

func (c *Controller) ListFavorites(ctx context.Context, uid uuid.UUID) ([]*model.Favorite, error) {
	const op = "favorites.ListFavorites.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := make([]*model.Favorite, 0, 15)
	cacheKey := fmt.Sprintf(favoriteCacheKey, uid)
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetFavorites(ctx, uid)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find favorites", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to get favorites", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err))
		}
	}

	return res, nil
}

func (c *Controller) AddToFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) (*model.Favorite, error) {
	const op = "favorites.AddToFavorites.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.AddToFavorites(ctx, uid, itemID)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	} else if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find item", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to add to favorites", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if err := c.cache.Delete(ctx, fmt.Sprintf(favoriteCacheKey, uid)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	return res, nil
}

func (c *Controller) RemoveFromFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) error {
	const op = "favorites.RemoveFromFavorites.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.RemoveFromFavorites(ctx, uid, itemID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find item", zap.Error(err), zap.String("op", op))
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to remove from favorites", zap.Error(err), zap.String("op", op))
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(favoriteCacheKey, uid)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	return nil
}
