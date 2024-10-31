package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	md "github.com/JMURv/par-pro/products/pkg/model"
	dbutils "github.com/JMURv/par-pro/products/pkg/utils/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"strings"
)

func (r *Repository) ListCategoryItems(ctx context.Context, slug string, page, size int, filters map[string]any, sort string) (*md.PaginatedItemsData, error) {
	const op = "items.ListCategoryItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	q := &strings.Builder{}
	q.WriteString(
		`
		SELECT COUNT(*)
		FROM item i
		JOIN item_category ic ON ic.item_id = i.id
		JOIN category c ON c.slug = ic.category_slug
		JOIN item_attr ia ON i.id = ia.item_id
		WHERE ic.category_slug = ?
	`,
	)

	args := make([]any, 0, 20)
	args = append(args, slug)
	args = dbutils.FilterItems(q, args, filters)

	var count int64
	if err := r.conn.QueryRow(q.String(), args...).Scan(&count); err != nil {
		return nil, err
	}
	q.Reset()
	args = []any{}

	q.WriteString(
		`
		SELECT i.id, i.title, i.price, i.src, i.alt, 
		ARRAY_AGG(c.title || '|' || c.slug) AS categories,
		ARRAY_AGG(ia.name || '|' || ia.value) AS attrs
		FROM item i
		JOIN item_category ic ON ic.item_id = i.id
		JOIN category c ON c.slug = ic.category_slug
		JOIN item_attr ia ON i.id = ia.item_id
		WHERE ic.category_slug = ?
	`,
	)

	args = append(args, slug)
	args = dbutils.FilterItems(q, args, filters)

	if sort == "" {
		sort = "i.created_at DESC"
	}

	q.WriteString(" GROUP BY i.id, i.created_at")
	q.WriteString(" ORDER BY ")
	q.WriteString(sort)

	q.WriteString(" OFFSET ? LIMIT ?;")
	args = append(args, (page-1)*size, size)

	rows, err := r.conn.Query(q.String(), args...)
	if err != nil {
		return nil, err
	}

	res := make([]*md.Item, 0, size)
	for rows.Next() {
		item := &md.Item{}
		categories := make([]string, 0, 10)
		attrs := make([]string, 0, 10)
		if err = rows.Scan(
			&item.ID,
			&item.Title,
			&item.Price,
			&item.Src,
			&item.Alt,
			pq.Array(&categories),
			pq.Array(&attrs),
		); err != nil {
			return nil, err
		}

		item.Categories, err = ScanItemCategories(categories)
		if err != nil {
			return nil, err
		}

		item.Attributes, err = ScanAttrs(attrs)
		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &md.PaginatedItemsData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ItemAttrSearch(ctx context.Context, query string, size, page int) (*md.PaginatedItemAttrData, error) {
	const op = "items.ItemAttrSearch.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(itemCountAttrsQ, "%"+query+"%").Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(itemSearchAttrQ, "%"+query+"%", (page-1)*size, size)
	if err != nil {
		return nil, err
	}

	attrs := make([]*md.ItemAttribute, 0, size)
	for rows.Next() {
		attr := &md.ItemAttribute{}
		if err = rows.Scan(
			&attr.ID,
			&attr.Name,
			&attr.Value,
		); err != nil {
			return nil, err
		}
		attrs = append(attrs, attr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
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

	terms := strings.Fields(query)
	conds := make([]string, 0, len(terms))
	args := make([]any, 0, len(terms)+2)
	for _, term := range terms {
		conds = append(conds, "title ILIKE ?")
		args = append(args, "%"+term+"%")
	}
	searchQ := strings.Join(conds, " AND ")

	var q strings.Builder
	q.WriteString("SELECT COUNT(*) FROM item WHERE ")
	q.WriteString(searchQ)

	var count int64
	if err := r.conn.QueryRow(q.String(), args...).Scan(&count); err != nil {
		return nil, err
	}
	q.Reset()

	q.WriteString(
		`
	SELECT 
	    i.id,
	    i.title,
	    i.article,
	    i.description,
	    i.price,
	    i.src,
	    i.alt,
	    i.in_stock,
	    i.is_hit,
	    i.is_rec,
	    i.parent_id,
	    i.created_at,
	    i.updated_at,
	    ARRAY_AGG(c.title || '|' || c.slug) AS categories
		FROM item i
		JOIN item_category ic ON i.id = ic.item_id
		JOIN category c ON c.slug = ic.category_slug
		WHERE 
   `,
	)
	q.WriteString(searchQ)
	q.WriteString(" GROUP BY i.id")
	q.WriteString(" OFFSET ? LIMIT ?")

	args = append(args, (page-1)*size, size)
	rows, err := r.conn.Query(q.String(), args...)
	if err != nil {
		return nil, err
	}

	items := make([]*md.Item, 0, size)
	for rows.Next() {
		item := &md.Item{}
		categories := make([]string, 0, 10)
		if err = rows.Scan(
			&item.ID,
			&item.Title,
			&item.Article,
			&item.Description,
			&item.Price,
			&item.Src,
			&item.Alt,
			&item.InStock,
			&item.IsHit,
			&item.IsRec,
			&item.ParentItemID,
			&item.CreatedAt,
			&item.UpdatedAt,
			pq.Array(&categories),
		); err != nil {
			return nil, err
		}

		item.Categories, err = ScanItemCategories(categories)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
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
	if err := r.conn.QueryRow(itemCountQ).Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(itemListQ, (page-1)*size, size)
	if err != nil {
		return nil, err
	}

	res := make([]*md.Item, 0, size)
	for rows.Next() {
		item := &md.Item{}
		categories := make([]string, 0, 10)
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Article,
			&item.Description,
			&item.Price,
			&item.Src,
			&item.Alt,
			&item.InStock,
			&item.IsHit,
			&item.IsRec,
			&item.ParentItemID,
			&item.CreatedAt,
			&item.UpdatedAt,
			pq.Array(&categories),
		); err != nil {
			return nil, err
		}

		item.Categories, err = ScanItemCategories(categories)
		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
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
	var err error
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	row := r.conn.QueryRow(itemGetByIDQ, uid)
	if err = row.Err(); err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	res := &md.Item{}
	media := make([]string, 0, 10)
	attrs := make([]string, 0, 10)
	categories := make([]string, 0, 10)
	if err := row.Scan(
		&res.ID,
		&res.Title,
		&res.Article,
		&res.Description,
		&res.Price,
		&res.Src,
		&res.Alt,
		&res.InStock,
		&res.IsHit,
		&res.IsRec,
		&res.ParentItemID,
		&res.CreatedAt,
		&res.UpdatedAt,
		pq.Array(&media),
		pq.Array(&attrs),
		pq.Array(&categories),
	); err != nil {
		return nil, err
	}

	res.Media, err = ScanMedia(media)
	if err != nil {
		return nil, err
	}

	res.Attributes, err = ScanAttrs(attrs)
	if err != nil {
		return nil, err
	}

	res.Categories, err = ScanItemCategories(categories)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) CreateItem(ctx context.Context, i *md.Item) (uuid.UUID, error) {
	const op = "items.CreateItem.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	if err = tx.QueryRow(
		itemCreateQ,
		i.Title,
		i.Article,
		i.Description,
		i.Price,
		i.Src,
		i.Alt,
		i.InStock,
		i.IsHit,
		i.IsRec,
		i.ParentItemID,
	).Scan(&id); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err = CreateItemMedia(tx, id, i.Media); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err = CreateItemAttrs(tx, id, i.Attributes); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err = CreateItemCategories(tx, id, i.Categories); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if len(i.Variants) > 0 {
		for _, v := range i.Variants {
			if _, err = r.conn.Exec(
				itemCreateQ,
				v.Title,
				v.Article,
				v.Description,
				v.Price,
				v.Src,
				v.Alt,
				v.InStock,
				v.IsHit,
				v.IsRec,
				id,
			); err != nil {
				tx.Rollback()
				return uuid.Nil, err
			}
		}
	}

	if len(i.RelatedProducts) > 0 {
		for _, v := range i.RelatedProducts {
			if _, err = r.conn.Exec(itemRelatedProductCreateQ, id, v.RelatedItemID); err != nil {
				tx.Rollback()
				return uuid.Nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return id, nil
}

func (r *Repository) UpdateItem(ctx context.Context, uid uuid.UUID, req *md.Item) error {
	const op = "items.UpdateItem.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		itemUpdateQ,
		req.Title,
		req.Article,
		req.Description,
		req.Price,
		req.Src,
		req.Alt,
		req.InStock,
		req.IsHit,
		req.IsRec,
		req.ParentItemID,
		uid,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		tx.Rollback()
		return repo.ErrNotFound
	}

	if err = UpdateItemMedia(tx, uid, req.Media); err != nil {
		tx.Rollback()
		return err
	}

	if err = UpdateItemAttributes(tx, uid, req.Attributes); err != nil {
		tx.Rollback()
		return err
	}

	if err = UpdateItemCategories(tx, uid, req.Categories); err != nil {
		tx.Rollback()
		return err
	}

	if err = UpdateItemRelatedProducts(tx, uid, req.RelatedProducts); err != nil {
		tx.Rollback()
		return err
	}

	if err = UpdateItemVariants(tx, uid, req.Variants); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) DeleteItem(ctx context.Context, uid uuid.UUID) error {
	const op = "items.DeleteItem.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.Exec(itemDeleteQ, uid)
	if err != nil {
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		return repo.ErrNotFound
	}

	return nil
}

func (r *Repository) ListItemVariants(ctx context.Context, uid uuid.UUID) ([]*md.Item, error) {
	const op = "items.ListItemVariants.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	rows, err := r.conn.Query(itemListVarsQ, uid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*md.Item, 0, 15)
	for rows.Next() {
		var i md.Item
		media := make([]string, 0, 10)
		attrs := make([]string, 0, 10)
		if err = rows.Scan(
			&i.ID,
			&i.Title,
			&i.Article,
			&i.Description,
			&i.Price,
			&i.Src,
			&i.Alt,
			&i.InStock,
			&i.IsHit,
			&i.IsRec,
			&i.ParentItemID,
			pq.Array(&media),
			pq.Array(&attrs),
		); err != nil {
			return nil, err
		}

		i.Media, err = ScanMedia(media)
		if err != nil {
			return nil, err
		}

		i.Attributes, err = ScanAttrs(attrs)
		if err != nil {
			return nil, err
		}

		res = append(res, &i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) ListRelatedItems(ctx context.Context, uid uuid.UUID) ([]*md.RelatedProduct, error) {
	const op = "items.ListRelatedItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	rows, err := r.conn.Query(itemListRelated, uid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]*md.RelatedProduct, 0, 15)
	for rows.Next() {
		var rp md.RelatedProduct
		if err = rows.Scan(
			&rp.RelatedItem.ID,
			&rp.RelatedItem.Title,
			&rp.RelatedItem.Price,
			&rp.RelatedItem.Src,
			&rp.RelatedItem.Alt,
		); err != nil {
			return nil, err
		}
		res = append(res, &rp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) HitItems(ctx context.Context, page, size int) (*md.PaginatedItemsData, error) {
	const op = "items.HitItems.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(itemCountHitQ).Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(itemListHitQ, (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*md.Item, 0, size)
	for rows.Next() {
		var i md.Item
		if err = rows.Scan(
			&i.ID,
			&i.Title,
			&i.Price,
			&i.Src,
			&i.Alt,
			&i.IsHit,
			&i.IsRec,
		); err != nil {
			return nil, err
		}
		res = append(res, &i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
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
	if err := r.conn.QueryRow(itemCountRecQ).Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(itemListRecQ, (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*md.Item, 0, size)
	for rows.Next() {
		var i md.Item
		if err = rows.Scan(
			&i.ID,
			&i.Title,
			&i.Price,
			&i.Src,
			&i.Alt,
			&i.IsHit,
			&i.IsRec,
		); err != nil {
			return nil, err
		}
		res = append(res, &i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &md.PaginatedItemsData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}
