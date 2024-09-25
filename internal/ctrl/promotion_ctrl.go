package ctrl

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
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
	ListPromotions(ctx context.Context, page, size int) (*utils.PaginatedData, error)
	GetPromotion(ctx context.Context, slug string) (*model.Promotion, error)
	CreatePromotion(ctx context.Context, p *model.Promotion) (*model.Promotion, error)
	UpdatePromotion(ctx context.Context, slug string, p *model.Promotion) (*model.Promotion, error)
	DeletePromotion(ctx context.Context, slug string) error

	ListPromotionItems(ctx context.Context, slug string, page, size int) (*utils.PaginatedData, error)

	PromotionSearch(ctx context.Context, query string, page int, size int) (*utils.PaginatedData, error)
}

func (c *Controller) invalidatePromoRelatedCache() {
	ctx := context.Background()
	if err := c.cache.InvalidateKeysByPattern(ctx, invalidatePromoRelatedCachePattern); err != nil {
		zap.L().Debug("failed to invalidate cache", zap.String("key", invalidatePromoRelatedCachePattern), zap.Error(err))
	}
}

func (c *Controller) PromotionSearch(ctx context.Context, query string, page int, size int) (*utils.PaginatedData, error) {
	const op = "promotion.search.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &utils.PaginatedData{}
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

func (c *Controller) ListPromotionItems(ctx context.Context, slug string, page, size int) (*utils.PaginatedData, error) {
	const op = "promo.ListPromotionItems.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &utils.PaginatedData{}
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

func (c *Controller) ListPromotions(ctx context.Context, page, size int) (*utils.PaginatedData, error) {
	const op = "promo.ListPromotions.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &utils.PaginatedData{}
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

func (c *Controller) CreatePromotion(ctx context.Context, p *model.Promotion) (*model.Promotion, error) {
	const op = "promo.CreatePromotion.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	p.Slug = slugify.Slugify(p.Title)

	// TODO: call mcrsvc
	p.Banner.PromotionSlug = &p.Slug
	p.SEO.PromotionSlug = &p.Slug
	res, err := c.repo.CreatePromotion(ctx, p)
	if err != nil {
		zap.L().Debug("failed to create promotion", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(promotionCacheKey, p.Slug), bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	go c.invalidatePromoRelatedCache()
	return res, nil
}

func (c *Controller) UpdatePromotion(ctx context.Context, slug string, p *model.Promotion) (*model.Promotion, error) {
	const op = "promo.UpdatePromotion.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.UpdatePromotion(ctx, slug, p)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find promotion", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to update promotion", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	// TODO: call mcrsvc
	banner, err := c.UpdateBannerByPromoSlug(ctx, res.Slug, &p.Banner)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to update promotion banner", zap.Error(err), zap.String("op", op))
		return nil, err
	}
	res.Banner = *banner

	// TODO: call mcrsvc
	seo, err := c.UpdatePromotionSEO(ctx, res.Slug, &p.SEO)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to update promotion SEO", zap.Error(err), zap.String("op", op))
		return nil, err
	}
	res.SEO = *seo

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(promotionCacheKey, res.Slug), bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	go c.invalidatePromoRelatedCache()
	return res, nil
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

	go c.invalidatePromoRelatedCache()
	return nil
}
