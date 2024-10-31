package ctrl

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/ctrl/seo"
	repo "github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/consts"
	md "github.com/JMURv/par-pro/products/pkg/model"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const itemCacheKey = "item:%v"
const itemSearchCacheKey = "items-search:%v:%v:%v"
const itemListCacheKey = "items-list:%v:%v"
const relatedItemCacheKey = "items-related:%v"
const recKey = "items-rec"
const hitKey = "items-hit"
const itemAttrSearchCacheKey = "items-attr-search:%v:%v:%v"
const itemCategoryCacheKey = "items-category:%v:%v:%v:%v:%v"
const invalidateItemRelatedCachePattern = "items-*"

type itemRepo interface {
	ListCategoryItems(ctx context.Context, slug string, page, size int, filters map[string]any, sort string) (*md.PaginatedItemsData, error)
	ItemAttrSearch(ctx context.Context, query string, size, page int) (res *md.PaginatedItemAttrData, err error)
	ItemSearch(ctx context.Context, query string, size, page int) (*md.PaginatedItemsData, error)

	ListItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error)
	GetItemByUUID(ctx context.Context, uid uuid.UUID) (*md.Item, error)
	CreateItem(ctx context.Context, i *md.Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, uid uuid.UUID, i *md.Item) error
	DeleteItem(ctx context.Context, uid uuid.UUID) error

	ListRelatedItems(ctx context.Context, uid uuid.UUID) ([]*md.RelatedProduct, error)

	HitItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error)
	RecItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error)
}

func (c *Controller) invalidateItemRelatedCache() {
	ctx := context.Background()
	if err := c.cache.InvalidateKeysByPattern(ctx, invalidateItemRelatedCachePattern); err != nil {
		zap.L().Debug(
			"failed to invalidate cache",
			zap.String("key", invalidateItemRelatedCachePattern),
			zap.Error(err),
		)
	}
}

func (c *Controller) ListCategoryItems(ctx context.Context, slug string, page, size int, filters map[string]any, sort string) (*md.PaginatedItemsData, error) {
	const op = "items.ListCategoryItems.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &md.PaginatedItemsData{}
	cacheKey := fmt.Sprintf(itemCategoryCacheKey, slug, page, size, filters, sort)
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ListCategoryItems(ctx, slug, page, size, filters, sort)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find category", zap.Error(err), zap.String("op", op))
		return nil, err
	} else if err != nil {
		zap.L().Debug("failed to list category items", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) ItemAttrSearch(ctx context.Context, query string, size, page int) (*md.PaginatedItemAttrData, error) {
	const op = "items.ItemAttrSearch.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &md.PaginatedItemAttrData{}
	if err := c.cache.GetToStruct(ctx, fmt.Sprintf(itemAttrSearchCacheKey, query, page, size), &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ItemAttrSearch(ctx, query, size, page)
	if err != nil {
		zap.L().Debug("failed to search items attributes", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(
			ctx,
			consts.DefaultCacheTime,
			fmt.Sprintf(itemAttrSearchCacheKey, query, page, size),
			bytes,
		); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) ItemSearch(ctx context.Context, query string, size, page int) (*md.PaginatedItemsData, error) {
	const op = "items.Search.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &md.PaginatedItemsData{}
	if err := c.cache.GetToStruct(ctx, fmt.Sprintf(itemSearchCacheKey, query, page, size), &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ItemSearch(ctx, query, page, size)
	if err != nil {
		zap.L().Debug("failed to search items", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(
			ctx,
			consts.DefaultCacheTime,
			fmt.Sprintf(itemSearchCacheKey, query, page, size),
			bytes,
		); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) ListRelatedItems(ctx context.Context, uid uuid.UUID) ([]*md.RelatedProduct, error) {
	const op = "items.ListRelatedItems.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := make([]*md.RelatedProduct, 0, 15)
	if err := c.cache.GetToStruct(ctx, fmt.Sprintf(relatedItemCacheKey, uid), &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ListRelatedItems(ctx, uid)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to find item", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to get related items", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(relatedItemCacheKey, uid), bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) HitItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.HitItems.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &md.PaginatedItemsData{}
	if err := c.cache.GetToStruct(ctx, hitKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.HitItems(ctx, page, size)
	if err != nil {
		zap.L().Debug("failed to get hit items", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, hitKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil

}

func (c *Controller) RecItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.RecItems.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &md.PaginatedItemsData{}
	if err := c.cache.GetToStruct(ctx, recKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.RecItems(ctx, page, size)
	if err != nil {
		zap.L().Debug("failed to get recommended items", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, recKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil

}

func (c *Controller) ListItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.ListItems.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &md.PaginatedItemsData{}
	if err := c.cache.GetToStruct(ctx, fmt.Sprintf(itemListCacheKey, page, size), cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ListItems(ctx, page, size)
	if err != nil {
		zap.L().Debug("failed to list items", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(
			ctx,
			consts.DefaultCacheTime,
			fmt.Sprintf(itemListCacheKey, page, size),
			bytes,
		); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) GetItemByUUID(ctx context.Context, uid uuid.UUID) (*md.Item, error) {
	const op = "items.GetItemByUUID.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &md.Item{}
	if err := c.cache.GetToStruct(ctx, fmt.Sprintf(itemCacheKey, uid), cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetItemByUUID(ctx, uid)
	if err != nil && err == repo.ErrNotFound {
		zap.L().Debug("failed to find item", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to get item", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(itemCacheKey, uid), bytes); err != nil {
			zap.L().Debug("failed to set cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) CreateItem(ctx context.Context, i *md.Item) (uuid.UUID, error) {
	const op = "items.CreateItem.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	uid, err := c.repo.CreateItem(ctx, i)
	if err != nil {
		zap.L().Debug("failed to create item", zap.Error(err), zap.String("op", op))
		return uuid.Nil, err
	}

	if err := c.seo.CreateSEO(ctx, seo.Item.String(), uid.String(), &i.SEO); err != nil {
		zap.L().Debug("failed to create item SEO", zap.Error(err), zap.String("op", op))
	}

	go c.invalidateItemRelatedCache()
	return uid, nil
}

func (c *Controller) UpdateItem(ctx context.Context, uid uuid.UUID, i *md.Item) error {
	const op = "items.UpdateItem.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.UpdateItem(ctx, uid, i)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to update item", zap.Error(err), zap.String("op", op))
		return err
	}

	if err := c.seo.UpdateSEO(ctx, seo.Item.String(), uid.String(), &i.SEO); err != nil {
		zap.L().Debug("failed to update item SEO", zap.Error(err), zap.String("op", op))
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(itemCacheKey, uid)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	go c.invalidateItemRelatedCache()
	return nil
}

func (c *Controller) DeleteItem(ctx context.Context, uid uuid.UUID) error {
	const op = "items.DeleteItem.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.DeleteItem(ctx, uid)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to delete item", zap.Error(err), zap.String("op", op))
		return err
	}

	if err := c.seo.DeleteSEO(ctx, seo.Item.String(), uid.String()); err != nil {
		zap.L().Debug("failed to delete item SEO", zap.Error(err), zap.String("op", op))
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(itemCacheKey, uid)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	go c.invalidateItemRelatedCache()
	return nil
}
