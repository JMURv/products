package db

import (
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) GetFavorites(ctx context.Context, uid uuid.UUID) ([]*model.Favorite, error) {
	const op = "favorites.GetFavorites.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var res []*model.Favorite
	if err := r.conn.Preload("Item").Where("user_id=?", uid).Find(&res).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) AddToFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) (*model.Favorite, error) {
	var err error
	const op = "favorites.AddToFavorites.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Favorite{}
	if err = r.conn.Where("user_id=? AND item_id=?", uid, itemID).First(res).Error; err == nil {
		return nil, repo.ErrAlreadyExists
	}

	item := &model.Item{}
	if err = r.conn.Where("id = ?", itemID).First(item).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	res = &model.Favorite{
		UserID:    uid,
		ItemID:    itemID,
		Item:      *item,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err = r.conn.Create(res).Error; err != nil {
		return nil, err
	}

	return res, nil

}

func (r *Repository) RemoveFromFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) error {
	const op = "favorites.RemoveFromFavorites.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Favorite{}
	if err := r.conn.Where("user_id=? AND item_id=?", uid, itemID).First(res).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	if err := r.conn.Delete(res).Error; err != nil {
		return err
	}

	return nil
}
