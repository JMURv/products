package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) PromotionSearch(ctx context.Context, query string, page, size int) (*model.PaginatedPromosData, error) {
	const op = "promo.search.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(promoSearchCountQ, "%"+query+"%", "%"+query+"%").Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(promoSearchQ, "%"+query+"%", "%"+query+"%", (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.Promotion, 0, size)
	for rows.Next() {
		p := &model.Promotion{}
		if err = rows.Scan(
			&p.Slug,
			&p.Title,
			&p.Description,
			&p.Src,
			&p.Alt,
			&p.LastsTo,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedPromosData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ListPromotionItems(ctx context.Context, slug string, page, size int) (*model.PaginatedPromoItemsData, error) {
	const op = "promo.ListPromotionItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(promoCountItemsQ, slug).Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(promoItemsQ, slug, (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.PromotionItem, 0, size)
	for rows.Next() {
		pi := &model.PromotionItem{}
		if err = rows.Scan(
			&pi.Discount,
			&pi.PromotionSlug,
			&pi.ItemID,
			&pi.Item.Title,
			&pi.Item.Price,
			&pi.Item.Src,
			&pi.Item.Alt,
		); err != nil {
			return nil, err
		}

		res = append(res, pi)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedPromoItemsData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) GetPromotion(ctx context.Context, slug string) (*model.Promotion, error) {
	const op = "promo.GetPromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Promotion{}
	err := r.conn.QueryRow(promoGetQ, slug).Scan(
		&res.Slug,
		&res.Title,
		&res.Description,
		&res.Src,
		&res.Alt,
		&res.LastsTo,
	)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) ListPromotions(ctx context.Context, page, size int) (*model.PaginatedPromosData, error) {
	const op = "promo.ListPromotions.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(promoCountQ).Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(promoListQ, (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.Promotion, 0, size)
	for rows.Next() {
		p := &model.Promotion{}
		if err = rows.Scan(
			&p.Slug,
			&p.Title,
			&p.Description,
			&p.Src,
			&p.Alt,
			&p.LastsTo,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedPromosData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) CreatePromotion(ctx context.Context, req *model.Promotion) (string, error) {
	const op = "promo.CreatePromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return "", err
	}

	var slug string
	err = tx.QueryRow(promoCreateQ, req.Slug, req.Title, req.Description, req.Src, req.Alt, req.LastsTo).Scan(&slug)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err = CreatePromotionItems(tx, slug, req.PromotionItems); err != nil {
		tx.Rollback()
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}
	return slug, nil
}

func (r *Repository) UpdatePromotion(ctx context.Context, slug string, req *model.Promotion) error {
	const op = "promo.UpdatePromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		promoUpdateQ,
		req.Title,
		req.Description,
		req.Src, req.Alt,
		req.LastsTo,
		slug,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		tx.Rollback()
		return repo.ErrNotFound
	}

	if err = UpdatePromotionItems(tx, slug, req.PromotionItems); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) DeletePromotion(ctx context.Context, slug string) error {
	const op = "promo.DeletePromotion.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, promoDeleteQ, slug)
	if err != nil {
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		return repo.ErrNotFound
	}

	return nil
}
