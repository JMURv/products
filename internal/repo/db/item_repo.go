package db

import (
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	md "github.com/JMURv/par-pro/products/pkg/model"
	gormutil "github.com/JMURv/par-pro/products/pkg/utils/gorm"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) ListCategoryItems(ctx context.Context, slug string, page, size int, filters map[string]any, sort string) (*md.PaginatedItemsData, error) {
	const op = "items.ListCategoryItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	query := r.conn.
		Model(&md.Item{}).
		Joins("JOIN item_categories ON item_categories.item_id = items.id").
		Joins("JOIN categories ON categories.slug = item_categories.category_slug").
		Where("categories.slug = ?", slug).
		Preload("Attributes")

	query = gormutil.FilterItems(query, filters)
	if sort == "" {
		sort = "items.created_at DESC"
	}
	query.Order(sort)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))

	var res []*md.Item
	if err := query.Offset((page - 1) * size).Limit(size).Find(&res).Error; err != nil {
		return nil, err
	}

	return &md.PaginatedItemsData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ItemAttrSearch(ctx context.Context, query string, size, page int) (res *md.PaginatedItemAttrData, err error) {
	const op = "items.ItemAttrSearch.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err = r.conn.
		Model(&md.ItemAttribute{}).
		Where("name ILIKE ?", "%"+query+"%").
		Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))

	var attrs []*md.ItemAttribute
	if err = r.conn.
		Offset((page-1)*size).
		Limit(size).
		Where("name ILIKE ?", "%"+query+"%").
		Find(&attrs).Error; err != nil {
		return nil, err
	}

	return &md.PaginatedItemAttrData{
		Data:        attrs,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ItemSearch(ctx context.Context, query string, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.ItemSearch.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.Model(&md.Item{}).Where("title ILIKE ?", "%"+query+"%").Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))

	var items []*md.Item
	if err := r.conn.
		Offset((page-1)*size).
		Limit(size).
		Where("title ILIKE ?", "%"+query+"%").
		Preload("Categories").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return &md.PaginatedItemsData{
		Data:        items,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ListItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.ListItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.
		Model(&md.Item{}).
		Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))

	var res []*md.Item
	if err := r.conn.
		Offset((page - 1) * size).
		Limit(size).
		Where("src != ''").
		Preload("Categories").
		Order("created_at desc").
		Find(&res).
		Error; err != nil {
		return nil, err
	}

	return &md.PaginatedItemsData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) GetItemByUUID(ctx context.Context, uid uuid.UUID) (*md.Item, error) {
	const op = "items.GetItemByUUID.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &md.Item{}
	if err := r.conn.
		Preload("Media").
		Preload("Categories").
		Preload("Attributes").
		Preload("Variants").
		Preload("Variants.Attributes").
		Preload("Variants.Media").
		Preload("RelatedProducts").
		Preload("RelatedProducts.RelatedItem").
		Where("id=?", uid).First(res).Error; err != nil && err == gorm.ErrRecordNotFound {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) CreateItem(ctx context.Context, i *md.Item) (*md.Item, error) {
	const op = "items.CreateItem.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if err := r.conn.Create(i).Error; err != nil {
		return nil, err
	}

	return i, nil
}

func (r *Repository) UpdateItem(ctx context.Context, uid uuid.UUID, newData *md.Item) (*md.Item, error) {
	const op = "items.UpdateItem.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	i, err := r.GetItemByUUID(ctx, uid)
	if err != nil && err == repo.ErrNotFound {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if newData.Title != "" {
		i.Title = newData.Title
	}

	if newData.Description != "" {
		i.Description = newData.Description
	}

	if newData.Price != 0 {
		i.Price = newData.Price
	}

	if newData.Src != "" {
		i.Src = newData.Src
	}

	if newData.Alt != "" {
		i.Alt = newData.Alt
	}

	if newData.QuantityInStock != 0 {
		i.QuantityInStock = newData.QuantityInStock
	}

	newCategories := make(map[string]struct{}, len(newData.Categories))
	for _, v := range newData.Categories {
		newCategories[v.Slug] = struct{}{}
	}
	for _, category := range i.Categories {
		if _, found := newCategories[category.Slug]; !found {
			if err := r.conn.
				Table("item_categories").
				Where("category_slug=?", category.Slug).
				Delete("item_categories").Error; err != nil {
				zap.L().Debug("failed to delete slide", zap.String("op", op), zap.Error(err))
			}
		}
	}

	i.Categories = newData.Categories
	i.Media = newData.Media
	i.Attributes = newData.Attributes
	i.Variants = newData.Variants
	i.RelatedProducts = newData.RelatedProducts
	i.IsHit = newData.IsHit
	i.IsRec = newData.IsRec

	i.UpdatedAt = time.Now()
	if err = r.conn.Save(i).Error; err != nil {
		return nil, err
	}
	return i, nil
}

func (r *Repository) DeleteItem(ctx context.Context, uid uuid.UUID) error {
	const op = "items.DeleteItem.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &md.Item{}
	if err := r.conn.Where("id = ?", uid).First(res).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	if err := r.conn.Delete(res).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListRelatedItems(ctx context.Context, uid uuid.UUID) ([]*md.RelatedProduct, error) {
	const op = "items.ListRelatedItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var item md.Item
	if err := r.conn.First(&item, "id = ?", uid).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var res []*md.RelatedProduct
	if err := r.conn.
		Preload("RelatedItem").
		Where("item_id=?", uid).
		Find(&res).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) HitItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.HitItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.
		Model(&md.Item{}).
		Where("is_hit=?", true).
		Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))

	var res []*md.Item
	if err := r.conn.
		Offset((page-1)*size).
		Limit(size).
		Preload("Media").
		Where("is_hit=?", true).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return &md.PaginatedItemsData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) RecItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.RecItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.
		Model(&md.Item{}).
		Where("is_rec=?", true).
		Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))

	var res []*md.Item
	if err := r.conn.
		Offset((page-1)*size).
		Limit(size).
		Preload("Media").
		Where("is_rec=?", true).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return &md.PaginatedItemsData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}
