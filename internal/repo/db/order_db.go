package db

import (
	"context"
	"database/sql"
	"errors"
	repo "github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	dbutils "github.com/JMURv/par-pro/products/pkg/utils/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"strings"
)

func (r *Repository) ListOrders(ctx context.Context, page, size int, filters map[string]any, sort string) (*model.PaginatedOrderData, error) {
	const op = "orders.ListOrders.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	q := &strings.Builder{}
	q.WriteString(
		`SELECT COUNT(*) FROM "order" o JOIN "order_item" oi ON o.id = oi.order_id JOIN item i ON oi.item_id = i.id WHERE 1=1`,
	)
	args := make([]any, 0, 20)
	args = dbutils.FilterOrders(q, args, filters)

	var count int64
	if err := r.conn.QueryRow(q.String(), args...).Scan(&count); err != nil {
		return nil, err
	}
	q.Reset()
	args = []any{}

	q.WriteString(
		`
		SELECT 
		    o.id, 
		    o.created_at, 
		    o.status,
		    ARRAY_AGG(oi.id || '|' || oi.quantity) AS order_items, 
		    ARRAY_AGG(i.id || '|' || i.title) AS items 
		FROM "order" o 
		JOIN "order_item" oi ON o.id = oi.order_id
		JOIN item i ON oi.item_id = i.id 
		WHERE 1=1 `,
	)
	args = dbutils.FilterOrders(q, args, filters)
	if sort == "" {
		sort = "o.created_at DESC"
	}
	q.WriteString(" GROUP BY o.id, o.created_at, o.status")
	q.WriteString(" ORDER BY ")
	q.WriteString(sort)
	q.WriteString(" OFFSET ? LIMIT ?;")
	args = append(args, (page-1)*size, size)

	rows, err := r.conn.Query(q.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.Order, 0, size)
	for rows.Next() {
		order := &model.Order{}
		orderItems := make([]string, 0, 10)
		items := make([]string, 0, 10)
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.Status,
			pq.Array(&orderItems),
			pq.Array(&items),
		); err != nil {
			return nil, err
		}

		order.OrderItems, err = ScanOrderItems(orderItems)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(order.OrderItems); i++ {
			order.OrderItems[i].Item, err = PopulateOrderItems(r.conn, order.OrderItems[i].ItemID)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedOrderData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) ListUserOrders(ctx context.Context, uid uuid.UUID, page, size int) (*model.PaginatedOrderData, error) {
	const op = "orders.ListUserOrders.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var count int64
	err := r.conn.QueryRow(userOrderCountQ, uid).Scan(&count)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	res := make([]*model.Order, 0, size)
	rows, err := r.conn.Query(userOrdersQ, uid, size, (page-1)*size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &model.Order{}
		items := make([]string, 0, 10)
		if err = rows.Scan(
			&order.ID,
			&order.Status,
			&order.TotalAmount,
			&order.FIO,
			&order.Tel,
			&order.Email,
			&order.Address,
			&order.Delivery,
			&order.PaymentMethod,
			&order.UserID,
			&order.CreatedAt,
			&order.UpdatedAt,
			pq.Array(&items),
		); err != nil {
			return nil, err
		}

		order.OrderItems, err = ScanOrderItems(items)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(order.OrderItems); i++ {
			order.OrderItems[i].Item, err = PopulateOrderItems(r.conn, order.OrderItems[i].ItemID)
			if err != nil {
				return nil, err
			}
		}

		res = append(res, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedOrderData{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) GetOrder(ctx context.Context, orderID uint64) (*model.Order, error) {
	const op = "orders.GetOrder.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Order{}
	items := make([]string, 0, 10)
	err := r.conn.QueryRow(orderGetQ, orderID).Scan(
		&res.ID,
		&res.Status,
		&res.TotalAmount,
		&res.FIO,
		&res.Tel,
		&res.Email,
		&res.Address,
		&res.Delivery,
		&res.PaymentMethod,
		&res.UserID,
		&res.CreatedAt,
		&res.UpdatedAt,
		pq.Array(&items),
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	res.OrderItems, err = ScanOrderItems(items)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(res.OrderItems); i++ {
		res.OrderItems[i].Item, err = PopulateOrderItems(r.conn, res.OrderItems[i].ItemID)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (r *Repository) CreateOrder(ctx context.Context, uid uuid.UUID, req *model.Order) (uint64, error) {
	const op = "orders.CreateOrder.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, v := range req.OrderItems {
		total += v.Item.Price * float64(v.Quantity)
	}

	orderID := uint64(0)
	err = tx.QueryRow(
		orderCreateQ,
		model.OrderStatusPending,
		total,
		req.FIO,
		req.Tel,
		req.Email,
		req.Address,
		req.Delivery,
		req.PaymentMethod,
		uid,
	).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for i := 0; i < len(req.OrderItems); i++ {
		if err = createOrderItem(tx, orderID, req.OrderItems[i]); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *Repository) UpdateOrder(ctx context.Context, orderID uint64, newData *model.Order) error {
	const op = "orders.UpdateOrder.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	if len(newData.OrderItems) > 0 {
		if err = UpdateOrderItems(tx, orderID, newData.OrderItems); err != nil {
			tx.Rollback()
			return err
		}

		var total float64
		orderItems := make([]model.OrderItem, len(newData.OrderItems))
		for i := 0; i < len(newData.OrderItems); i++ {
			v := newData.OrderItems[i]
			price := float64(0)

			if err = tx.QueryRow(orderItemPriceQ, v.ItemID).Scan(&price); err != nil && errors.Is(err, sql.ErrNoRows) {
				tx.Rollback()
				return repo.ErrNotFound
			} else if err != nil {
				tx.Rollback()
				return err
			}

			total += price * float64(v.Quantity)
			orderItems[i] = model.OrderItem{
				Quantity: v.Quantity,
				ItemID:   v.ItemID,
			}
		}
		newData.TotalAmount = total
	}

	res, err := tx.Exec(
		orderUpdateQ,
		newData.Status,
		newData.TotalAmount,
		newData.FIO,
		newData.Tel,
		newData.Email,
		newData.Address,
		newData.Delivery,
		newData.PaymentMethod,
		orderID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if eff, _ := res.RowsAffected(); eff == 0 {
		tx.Rollback()
		return repo.ErrNotFound
	}

	return tx.Commit()

}

func (r *Repository) CancelOrder(ctx context.Context, orderID uint64) error {
	const op = "orders.CancelOrder.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(orderCancelQ, model.OrderStatusCancelled, orderID)
	if err != nil {
		return err
	}

	if eff, _ := res.RowsAffected(); eff == 0 {
		return repo.ErrNotFound
	}

	return tx.Commit()
}
