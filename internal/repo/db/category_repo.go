package db

import (
	"context"
	"errors"
	"fmt"
	repo "github.com/JMURv/par-pro/internal/repository"
	"github.com/JMURv/par-pro/pkg/model"
	utils "github.com/JMURv/par-pro/pkg/utils/http"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) CategoryFiltersSearch(ctx context.Context, query string, page int, size int) (res *utils.PaginatedData, err error) {
	const op = "category.categoryFiltersSearch.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err = r.conn.
		Model(&model.Filter{}).
		Where("name ILIKE ?", "%"+query+"%").
		Count(&count).
		Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))

	var categories []*model.Filter
	if err = r.conn.
		Offset((page-1)*size).
		Limit(size).
		Where("name ILIKE ?", "%"+query+"%").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return &utils.PaginatedData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) CategorySearch(ctx context.Context, query string, page, size int) (*utils.PaginatedData, error) {
	const op = "category.search.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.
		Model(&model.Category{}).
		Where("title ILIKE ? OR slug ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).
		Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))
	hasNextPage := page < totalPages

	var res []*model.Category
	if err := r.conn.
		Offset((page-1)*size).
		Limit(size).
		Where("title ILIKE ? OR slug ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&res).Error; err != nil {
		return nil, err
	}

	return &utils.PaginatedData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: hasNextPage,
	}, nil
}

func (r *Repository) ListCategories(ctx context.Context, page, size int) (*utils.PaginatedData, error) {
	const op = "category.ListCategories.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.
		Model(&model.Category{}).
		Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))
	hasNextPage := page < totalPages

	var res []*model.Category
	if err := r.conn.
		Preload("Children").
		Offset((page - 1) * size).
		Limit(size).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return &utils.PaginatedData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: hasNextPage,
	}, nil
}

func (r *Repository) GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	const op = "category.GetCategoryBySlug.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Category{}
	if err := r.conn.
		Preload("Filters").
		Preload("Children").
		Preload("ParentCategory").
		Where("slug=?", slug).First(res).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	for i := range res.Children {
		childSlug := res.Children[i].Slug
		var count int64

		if err := r.conn.Model(&model.Item{}).
			Joins("JOIN item_categories ON item_categories.item_id = items.id").
			Where("item_categories.category_slug = ?", childSlug).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count items for category %s: %w", childSlug, err)
		}

		res.Children[i].ProductQuantity = int(count)
	}

	return res, nil
}

func (r *Repository) CreateCategory(ctx context.Context, c *model.Category) (*model.Category, error) {
	const op = "category.CreateCategory.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	if err := r.conn.Create(c).Error; err != nil {
		return nil, err
	}

	return c, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, slug string, newData *model.Category) (*model.Category, error) {
	const op = "category.UpdateCategory.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	c, err := r.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, repo.ErrNotFound
	}

	if newData.Slug != "" {
		c.Slug = newData.Slug
	}
	if newData.Title != "" {
		c.Title = newData.Title
	}
	if newData.ProductQuantity != 0 {
		c.ProductQuantity = newData.ProductQuantity
	}
	if newData.Src != "" {
		c.Src = newData.Src
	}
	if newData.Alt != "" {
		c.Alt = newData.Alt
	}

	if newData.ParentSlug != nil && *newData.ParentSlug != "" {
		newParent, err := r.GetCategoryBySlug(ctx, *newData.ParentSlug)
		if err != nil {
			return nil, repo.ErrNotFound
		}
		c.ParentCategory = newParent
	}

	if len(newData.Filters) != 0 {
		newFilters := make(map[uint64]struct{}, len(newData.Filters))
		for _, v := range newData.Filters {
			newFilters[v.ID] = struct{}{}
		}

		for _, v := range c.Filters {
			if _, found := newFilters[v.ID]; !found {
				if err := r.conn.Where("id=?", v.ID).Delete(&model.Filter{}).Error; err != nil {
					zap.L().Debug("failed to delete filter", zap.String("op", op), zap.Error(err))
				}
			}
		}

		c.Filters = newData.Filters
	}

	c.UpdatedAt = time.Now()
	if err = r.conn.Save(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repository) DeleteCategory(ctx context.Context, slug string) error {
	const op = "category.DeleteCategory.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Category{}
	if err := r.conn.Where("slug = ?", slug).First(res).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	if err := r.conn.Delete(res).Error; err != nil {
		return err
	}

	return nil
}

// Filters
// TODO: Высчитать тут макс и мин значения

func (r *Repository) ListCategoryFilters(ctx context.Context, slug string) ([]*model.Filter, error) {
	var filters []*model.Filter

	err := r.conn.Where("category_slug = ?", slug).Find(&filters).Error
	return filters, err
}
