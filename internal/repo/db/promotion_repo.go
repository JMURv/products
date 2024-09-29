package db

import (
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) PromotionSearch(ctx context.Context, query string, page, size int) (*utils.PaginatedData, error) {
	const op = "promo.search.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.
		Model(&model.Promotion{}).
		Where("title ILIKE ? OR slug ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).
		Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))
	hasNextPage := page < totalPages

	var res []*model.Promotion
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

func (r *Repository) ListPromotionItems(ctx context.Context, slug string, page, size int) (*utils.PaginatedData, error) {
	const op = "promo.ListPromotionItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.Model(&model.PromotionItem{}).Where("promotion_slug=?", slug).Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))
	hasNextPage := page < totalPages

	var res []*model.PromotionItem
	if err := r.conn.
		Offset((page-1)*size).
		Limit(size).
		Preload("Item").
		Where("promotion_slug=?", slug).
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

func (r *Repository) GetPromotion(ctx context.Context, slug string) (*model.Promotion, error) {
	const op = "promo.GetPromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Promotion{}
	if err := r.conn.
		Preload("PromotionItems").
		Preload("PromotionItems.Item").
		Where("slug=?", slug).
		First(res).Error; err != nil && err == gorm.ErrRecordNotFound {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) ListPromotions(ctx context.Context, page, size int) (*utils.PaginatedData, error) {
	const op = "promo.ListPromotions.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.Model(&model.Promotion{}).Count(&count).Error; err != nil {
		return nil, err
	}
	totalPages := int((count + int64(size) - 1) / int64(size))
	hasNextPage := page < totalPages

	var res []*model.Promotion
	if err := r.conn.
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

func (r *Repository) CreatePromotion(ctx context.Context, p *model.Promotion) (*model.Promotion, error) {
	const op = "promo.CreatePromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if err := r.conn.Create(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repository) UpdatePromotion(ctx context.Context, slug string, newData *model.Promotion) (*model.Promotion, error) {
	const op = "promo.UpdatePromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	p, err := r.GetPromotion(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx := r.conn.Begin()

	if newData.Slug != "" {
		p.Slug = newData.Slug
	}

	if newData.Title != "" {
		p.Title = newData.Title
	}

	if newData.Description != "" {
		p.Description = newData.Description
	}

	if newData.Src != "" {
		p.Src = newData.Src
	}

	if newData.Alt != "" {
		p.Alt = newData.Alt
	}

	if !newData.LastsTo.IsZero() {
		p.LastsTo = newData.LastsTo
	}

	newItems := make(map[uuid.UUID]struct{}, len(newData.PromotionItems))
	for _, v := range newData.PromotionItems {
		if v.ID == 0 {
			newItems[v.ItemID] = struct{}{}
		}
	}
	for _, v := range p.PromotionItems {
		if _, found := newItems[v.ItemID]; !found {
			if err := tx.Where("id=?", v.ID).Delete(&model.PromotionItem{}).Error; err != nil {
				zap.L().Debug("failed to delete item", zap.String("op", op), zap.Error(err))
			}
		}
	}

	p.PromotionItems = newData.PromotionItems
	p.UpdatedAt = time.Now()
	if err = tx.Save(p).Error; err != nil {
		return nil, err
	}

	tx.Commit()
	return p, nil
}

func (r *Repository) DeletePromotion(ctx context.Context, slug string) error {
	const op = "promo.DeletePromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Promotion{}
	if err := r.conn.Where("slug=?", slug).First(res).Error; err != nil && err == gorm.ErrRecordNotFound {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	if err := r.conn.Where("slug=?", slug).Delete(res).Error; err != nil {
		return err
	}

	return nil
}
