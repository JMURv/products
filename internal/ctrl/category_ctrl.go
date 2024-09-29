package ctrl

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/ctrl/etc"
	"github.com/JMURv/par-pro/products/internal/ctrl/seo"
	repo "github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/JMURv/par-pro/products/pkg/utils/slugify"
	"github.com/goccy/go-json"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const categoryCacheKey = "category:%v"
const categoryListCacheKey = "categories-list:%v:%v"
const categoryFiltersListCacheKey = "categories-filters-list:%v"
const categorySearchCacheKey = "categories-search:%v:%v:%v"
const categoryFiltersSearchCacheKey = "categories-filters-search:%v:%v:%v"
const invalidateCategoryRelatedCachePattern = "categories-*"

type categoryRepo interface {
	ListCategories(ctx context.Context, page, size int) (*utils.PaginatedData, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error)
	CreateCategory(ctx context.Context, c *model.Category) (*model.Category, error)
	UpdateCategory(ctx context.Context, slug string, c *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, slug string) error

	ListCategoryFilters(ctx context.Context, slug string) ([]*model.Filter, error)

	CategorySearch(ctx context.Context, query string, page int, size int) (*utils.PaginatedData, error)
	CategoryFiltersSearch(ctx context.Context, query string, page int, size int) (res *utils.PaginatedData, err error)
}

func (c *Controller) invalidateCategoryRelatedCache() {
	ctx := context.Background()
	if err := c.cache.InvalidateKeysByPattern(ctx, invalidateCategoryRelatedCachePattern); err != nil {
		zap.L().Debug("failed to invalidate cache", zap.String("key", invalidateCategoryRelatedCachePattern), zap.Error(err))
	}
}

func (c *Controller) CategoryFiltersSearch(ctx context.Context, query string, page int, size int) (*utils.PaginatedData, error) {
	const op = "category.categoryFiltersSearch.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cacheKey := fmt.Sprintf(categoryFiltersSearchCacheKey, query, page, size)
	cached := &utils.PaginatedData{}
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.CategoryFiltersSearch(ctx, query, page, size)
	if err != nil {
		zap.L().Debug("failed to search filters", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}
	return res, nil
}

func (c *Controller) CategorySearch(ctx context.Context, query string, page int, size int) (*utils.PaginatedData, error) {
	const op = "category.search.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cacheKey := fmt.Sprintf(categorySearchCacheKey, query, page, size)
	cached := &utils.PaginatedData{}
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.CategorySearch(ctx, query, page, size)
	if err != nil {
		zap.L().Debug("failed to search categories", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}
	return res, nil
}

func (c *Controller) ListCategoryFilters(ctx context.Context, slug string) ([]*model.Filter, error) {
	const op = "category.ListCategoryFilters.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cacheKey := fmt.Sprintf(categoryFiltersListCacheKey, slug)
	cached := make([]*model.Filter, 0, 15)
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ListCategoryFilters(ctx, slug)
	if err != nil {
		zap.L().Debug("failed to list category filters", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}
	return res, nil
}

func (c *Controller) ListCategories(ctx context.Context, page, size int) (*utils.PaginatedData, error) {
	const op = "category.ListCategories.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cacheKey := fmt.Sprintf(categoryListCacheKey, page, size)
	cached := &utils.PaginatedData{}
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ListCategories(ctx, page, size)
	if err != nil {
		zap.L().Debug("failed to list categories", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	const op = "category.GetCategoryBySlug.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.Category{}
	if err := c.cache.GetToStruct(ctx, fmt.Sprintf(categoryCacheKey, slug), cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetCategoryBySlug(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("failed to found category", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to get category", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(categoryCacheKey, slug), bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}

	return res, nil
}

func (c *Controller) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	const op = "category.CreateCategory.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	category.Slug = slugify.Slugify(category.Title)
	res, err := c.repo.CreateCategory(ctx, category)
	if err != nil {
		zap.L().Debug("failed to create category", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if err := c.etc.CreateBanner(ctx, etc.Category.String(), res.Slug, category.Banner); err != nil {
		zap.L().Debug("failed to create category banner", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if err := c.seo.CreateSEO(ctx, seo.Category.String(), res.Slug, category.SEO); err != nil {
		zap.L().Debug("failed to create category SEO", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(categoryCacheKey, res.Slug), bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err))
		}
	}

	if err := c.cache.Delete(ctx, fmt.Sprintf(categoryCacheKey, res.ParentSlug)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err))
	}

	go c.invalidateCategoryRelatedCache()
	return res, nil
}

func (c *Controller) UpdateCategory(ctx context.Context, slug string, category *model.Category) (*model.Category, error) {
	const op = "category.UpdateCategory.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.UpdateCategory(ctx, slug, category)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to update category", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if err := c.etc.UpdateBanner(ctx, etc.Category.String(), res.Slug, category.Banner); err != nil {
		zap.L().Debug("failed to update category banner", zap.Error(err), zap.String("op", op))
	}

	if err := c.seo.UpdateSEO(ctx, seo.Category.String(), res.Slug, res.SEO); err != nil {
		zap.L().Debug("failed to update category SEO", zap.Error(err), zap.String("op", op))
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(categoryCacheKey, res.Slug), bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err))
		}
	}

	go c.invalidateCategoryRelatedCache()
	return res, nil
}

func (c *Controller) DeleteCategory(ctx context.Context, slug string) error {
	const op = "category.DeleteCategory.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.DeleteCategory(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("failed to delete category", zap.Error(err), zap.String("op", op))
		return err
	}

	if err := c.etc.DeleteBanner(ctx, etc.Category.String(), slug); err != nil {
		zap.L().Debug("failed to delete category banner", zap.Error(err), zap.String("op", op))
	}

	if err := c.seo.DeleteSEO(ctx, seo.Category.String(), slug); err != nil {
		zap.L().Debug("failed to delete category SEO", zap.Error(err), zap.String("op", op))
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(categoryCacheKey, slug)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err))
	}

	go c.invalidateCategoryRelatedCache()
	return nil
}
