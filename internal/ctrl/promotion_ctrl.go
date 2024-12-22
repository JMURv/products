package ctrl

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/JMURv/par-pro/products/pkg/utils/slugify"
	"github.com/goccy/go-json"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const promotionCacheKey = "promo:%v"
const PromotionListCacheKey = "promos-list:%v:%v"
const promotionSearchCacheKey = "promos-search:%v:%v:%v"
const PromotionItemsCacheKey = "items-promos:%v:%v:%v"
const invalidatePromoRelatedCachePattern = "promos-*"

type promotionRepo interface {
	PromotionSearch(ctx context.Context, query string, page int, size int) (*model.PaginatedPromosData, error)
	ListPromotions(ctx context.Context, page, size int) (*model.PaginatedPromosData, error)
	GetPromotion(ctx context.Context, slug string) (*model.Promotion, error)
	CreatePromotion(ctx context.Context, p *model.Promotion) (string, error)
	UpdatePromotion(ctx context.Context, slug string, p *model.Promotion) error
	DeletePromotion(ctx context.Context, slug string) error

	ListPromotionItems(ctx context.Context, slug string, page, size int) (*model.PaginatedPromoItemsData, error)
}

func (c *Controller) PromotionSearch(ctx context.Context, query string, page int, size int) (*model.PaginatedPromosData, error) {
	const op = "promotion.search.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.PaginatedPromosData{}
	cacheKey := fmt.Sprintf(promotionSearchCacheKey, query, page, size)
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.PromotionSearch(ctx, query, page, size)
	if err != nil {
		zap.L().Debug("failed to search promotions", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) ListPromotions(ctx context.Context, page, size int) (*model.PaginatedPromosData, error) {
	const op = "promo.ListPromotions.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.PaginatedPromosData{}
	cacheKey := fmt.Sprintf(PromotionListCacheKey, page, size)
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ListPromotions(ctx, page, size)
	if err != nil {
		zap.L().Debug("failed to list promotions", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) GetPromotion(ctx context.Context, slug string) (*model.Promotion, error) {
	const op = "promo.GetPromotion.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.Promotion{}
	cacheKey := fmt.Sprintf(promotionCacheKey, slug)
	if err := c.cache.GetToStruct(ctx, cacheKey, cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetPromotion(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find promotion", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to get promotion", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) CreatePromotion(ctx context.Context, p *model.Promotion) (string, error) {
	const op = "promo.CreatePromotion.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	p.Slug = slugify.Slugify(p.Title)
	slug, err := c.repo.CreatePromotion(ctx, p)
	if err != nil {
		zap.L().Debug("failed to create promotion", zap.Error(err), zap.String("op", op))
		return "", err
	}

	go c.cache.InvalidateKeysByPattern(ctx, invalidatePromoRelatedCachePattern)
	return slug, nil
}

func (c *Controller) UpdatePromotion(ctx context.Context, slug string, p *model.Promotion) error {
	const op = "promo.UpdatePromotion.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.UpdatePromotion(ctx, slug, p)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find promotion", zap.Error(err), zap.String("op", op))
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to update promotion", zap.Error(err), zap.String("op", op))
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(promotionCacheKey, slug)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	go c.cache.InvalidateKeysByPattern(ctx, invalidatePromoRelatedCachePattern)
	return nil
}

func (c *Controller) DeletePromotion(ctx context.Context, slug string) error {
	const op = "promo.DeletePromotion.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.DeletePromotion(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find promotion", zap.Error(err), zap.String("op", op))
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to delete promotion", zap.Error(err), zap.String("op", op))
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(promotionCacheKey, slug)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	go c.cache.InvalidateKeysByPattern(ctx, invalidatePromoRelatedCachePattern)
	return nil
}

func (c *Controller) ListPromotionItems(ctx context.Context, slug string, page, size int) (*model.PaginatedPromoItemsData, error) {
	const op = "promo.ListPromotionItems.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.PaginatedPromoItemsData{}
	cacheKey := fmt.Sprintf(PromotionItemsCacheKey, slug, page, size)
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}
	res, err := c.repo.ListPromotionItems(ctx, slug, page, size)
	if err != nil {
		zap.L().Debug("failed to get promotion items", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}
	return res, nil
}
