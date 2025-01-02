package db

import (
	"database/sql"
	"errors"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
)

func createOrderItem(tx *sql.Tx, orderID uint64, req model.OrderItem) error {
	if _, err := tx.Exec(
		orderItemCreateQ,
		orderID,
		req.ItemID,
		req.Quantity,
	); err != nil {
		return err
	}

	return nil
}

func UpdateOrderItems(tx *sql.Tx, orderID uint64, req []model.OrderItem) error {
	rows, err := tx.Query(orderItemListQ, orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[uint64]*model.OrderItem)
	for rows.Next() {
		oi := &model.OrderItem{}
		if err = rows.Scan(&oi.ID, &oi.ItemID, &oi.Quantity); err != nil {
			return err
		}
		existing[oi.ID] = oi
	}

	if err = rows.Err(); err != nil {
		return err
	}

	reqMap := make(map[uint64]struct{}, len(req))
	for i := 0; i < len(req); i++ {
		reqMap[req[i].ID] = struct{}{}
		if _, ok := existing[req[i].ID]; !ok {
			if err = createOrderItem(tx, orderID, req[i]); err != nil {
				return err
			}
		}
	}

	for k, v := range existing {
		if _, ok := reqMap[k]; !ok {
			if _, err = tx.Exec(orderItemDeleteQ, v.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func PopulateOrderItems(tx *sql.DB, itemID uuid.UUID) (model.Item, error) {
	i := model.Item{}
	err := tx.QueryRow(
		`SELECT 
    		i.id, i.title, i.article, i.price, i.src, i.alt, i.in_stock
			FROM item i
			WHERE id = $1`,
		itemID,
	).Scan(
		&i.ID,
		&i.Title,
		&i.Article,
		&i.Price,
		&i.Src,
		&i.Alt,
		&i.InStock,
	)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return model.Item{}, repo.ErrNotFound
	}

	return i, nil
}
