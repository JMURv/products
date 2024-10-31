package db

import (
	"database/sql"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/lib/pq"
)

func createFilter(tx *sql.Tx, slug string, req model.Filter) error {
	if err := tx.QueryRow(
		filterCreateQ,
		req.Name,
		pq.Array(req.Values),
		req.FilterType,
		req.MinValue,
		req.MaxValue,
		slug,
	).Err(); err != nil {
		return err
	}

	return nil
}

func CreateFilters(tx *sql.Tx, slug string, req []model.Filter) error {
	for _, f := range req {
		if err := createFilter(tx, slug, f); err != nil {
			return err
		}
	}
	return nil
}

func UpdateFilters(tx *sql.Tx, slug string, req []model.Filter) error {
	rows, err := tx.Query(filterListQ, slug)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[uint64]model.Filter)
	for rows.Next() {
		var v model.Filter
		if err = rows.Scan(
			&v.ID,
			&v.Name,
			&v.Values,
			&v.FilterType,
			&v.MinValue,
			&v.MaxValue,
			&v.CategorySlug,
		); err != nil {
			return err
		}
		existing[v.ID] = v
	}

	if err = rows.Err(); err != nil {
		return err
	}

	for _, v := range req {
		if _, exists := existing[v.ID]; exists {
			if _, err = tx.Exec(
				filterUpdateQ,
				v.Name,
				pq.Array(v.Values),
				v.FilterType,
				v.MinValue,
				v.MaxValue,
				slug,
				v.ID,
			); err != nil {
				return err
			}
		} else {
			if err = createFilter(tx, slug, v); err != nil {
				return err
			}
		}
	}

	reqMap := make(map[uint64]struct{}, len(req))
	for _, v := range req {
		reqMap[v.ID] = struct{}{}
	}
	for id := range existing {
		if _, ok := reqMap[id]; !ok {
			if _, err = tx.Exec(filterDeleteQ, id, slug); err != nil {
				return err
			}
		}
	}

	return nil
}
