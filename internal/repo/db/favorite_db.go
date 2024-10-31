package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"strings"
)

func (r *Repository) GetFavorites(ctx context.Context, uid uuid.UUID) ([]*model.Favorite, error) {
	const op = "favorites.GetFavorites.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	rows, err := r.conn.Query(favGetQ, uid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*model.Favorite
	for rows.Next() {
		var f model.Favorite
		if err = rows.Scan(
			&f.UserID,
			&f.ItemID,
			&f.Item.ID,
			&f.Item.Title,
			&f.Item.Src,
			&f.Item.Alt,
			&f.Item.Price,
		); err != nil {
			return nil, err
		}
		res = append(res, &f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) AddToFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) (*model.Favorite, error) {
	var err error
	const op = "favorites.AddToFavorites.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var insertedItemID uuid.UUID
	err = r.conn.QueryRow(favAddQ, uid, itemID).Scan(&insertedItemID)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	if insertedItemID == uuid.Nil {
		return nil, repo.ErrAlreadyExists
	}

	return &model.Favorite{
		UserID: uid,
		ItemID: itemID,
	}, nil
}

func (r *Repository) RemoveFromFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) error {
	const op = "favorites.RemoveFromFavorites.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.Exec(favDelQ, uid, itemID)
	if err != nil {
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		return repo.ErrNotFound
	}

	return nil
}
