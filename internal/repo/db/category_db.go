package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"strings"
)

func (r *Repository) CategorySearch(ctx context.Context, query string, page, size int) (*model.PaginatedCategoryData, error) {
	const op = "category.search.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(categorySearchCountQ, "%"+query+"%", "%"+query+"%").Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(categorySearchQ, "%"+query+"%", "%"+query+"%", (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.Category, 0, size)
	for rows.Next() {
		c := &model.Category{}
		if err = rows.Scan(&c.Slug, &c.Title, &c.Src, &c.Alt, &c.ParentSlug); err != nil {
			return nil, err
		}
		res = append(res, c)
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedCategoryData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ListCategories(ctx context.Context, page, size int) (*model.PaginatedCategoryData, error) {
	const op = "category.ListCategories.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(categoryCountQ).Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(categoryListQ, (page-1)*size, size)
	if err != nil {
		return nil, err
	}

	res := make([]*model.Category, 0, size)
	for rows.Next() {
		c := &model.Category{}
		children := make([]string, 0, consts.DefaultPageSize)
		if err = rows.Scan(
			&c.Slug,
			&c.Title,
			&c.Src,
			&c.Alt,
			&c.ParentSlug,
			pq.Array(&children),
		); err != nil {
			return nil, err
		}

		c.Children, err = ScanChildrenCategory(children)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedCategoryData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	const op = "category.GetCategoryBySlug.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Category{}
	children := make([]string, 0, consts.DefaultPageSize)
	filters := make([]string, 0, consts.DefaultPageSize)
	err := r.conn.QueryRow(categoryGetQ, slug).
		Scan(&res.Slug, &res.Title, &res.Src, &res.Alt, &res.ParentSlug, pq.Array(&children), pq.Array(&filters))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	res.Children, err = ScanChildrenCategory(children)
	if err != nil {
		return nil, err
	}

	res.Filters, err = ScanFilters(filters)
	if err != nil {
		return nil, err
	}

	for i := range res.Children {
		var count int64
		err = r.conn.QueryRow(
			`SELECT COUNT(*) 
			 FROM item 
			 JOIN item_category ON item_category.item_id = item.id
			 WHERE item_category.category_slug = $1`, res.Children[i].Slug,
		).Scan(&count)
		if err != nil {
			return nil, err
		}

		res.Children[i].ProductQuantity = count
	}

	return res, nil
}

func (r *Repository) CreateCategory(ctx context.Context, req *model.Category) (string, error) {
	const op = "category.CreateCategory.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return "", err
	}

	var slug string
	if err = tx.QueryRow(categoryCreateQ, req.Slug, req.Title, req.Src, req.Alt, req.ParentSlug).
		Scan(&slug); err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "unique constraint") {
			return "", repo.ErrAlreadyExists
		}
		return "", err
	}

	if err = CreateFilters(tx, slug, req.Filters); err != nil {
		tx.Rollback()
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}
	return slug, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, slug string, req *model.Category) error {
	const op = "category.UpdateCategory.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	res, err := r.conn.Exec(categoryUpdateQ, req.Title, req.Src, req.Alt, req.ParentSlug, slug)
	if err != nil {
		tx.Rollback()
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		tx.Rollback()
		return repo.ErrNotFound
	}

	if err = UpdateFilters(tx, slug, req.Filters); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) DeleteCategory(ctx context.Context, slug string) error {
	const op = "category.DeleteCategory.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, categoryDeleteQ, slug)
	if err != nil {
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		return repo.ErrNotFound
	}

	return nil
}

// Filters
func (r *Repository) CategoryFiltersSearch(ctx context.Context, query string, page int, size int) (*model.PaginatedFilterData, error) {
	const op = "category.categoryFiltersSearch.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	if err := r.conn.QueryRow(filterCountQ, "%"+query+"%").Scan(&count); err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(filterSearchQ, "%"+query+"%", (page-1)*size, size)
	if err != nil {
		return nil, err
	}

	res := make([]*model.Filter, 0, size)
	for rows.Next() {
		f := &model.Filter{}
		if err = rows.Scan(
			&f.ID,
			&f.Name,
			&f.Values,
			&f.FilterType,
			&f.MinValue,
			&f.MaxValue,
			&f.CategorySlug,
		); err != nil {
			return nil, err
		}
		res = append(res, f)
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedFilterData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ListCategoryFilters(ctx context.Context, slug string) ([]*model.Filter, error) {
	const op = "category.ListCategoryFilters.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	rows, err := r.conn.Query(filterListQ, slug)
	if err != nil {
		return nil, err
	}

	res := make([]*model.Filter, 0, consts.DefaultPageSize)
	for rows.Next() {
		f := &model.Filter{}
		if err = rows.Scan(
			&f.ID,
			&f.Name,
			&f.Values,
			&f.FilterType,
			&f.MinValue,
			&f.MaxValue,
			&f.CategorySlug,
		); err != nil {
			return nil, err
		}
		res = append(res, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, err
}
