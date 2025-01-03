package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	repo2 "github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
	"time"
)

func TestRepository_ListUserOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()
	page := 1
	size := 10
	expectedCount := int64(2)
	expectedTotalPages := int((expectedCount + int64(size) - 1) / int64(size))

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.PaginatedOrderData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(userOrderCountQ)).
					WithArgs(uid).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id",
						"status",
						"total_amount",
						"fio",
						"tel",
						"email",
						"address",
						"delivery",
						"payment_method",
						"user_id",
						"created_at",
						"updated_at",
						"order_items",
					},
				).
					AddRow(
						1,
						"completed",
						100.0,
						"John Doe",
						"1234567890",
						"john@example.com",
						"123 Street",
						"delivery",
						"credit_card",
						uid.String(),
						time.Now(),
						time.Now(),
						fmt.Sprintf("{1|%v|3}", uid),
					).
					AddRow(
						2,
						"pending",
						50.0,
						"Jane Doe",
						"0987654321",
						"jane@example.com",
						"456 Street",
						"pickup",
						"cash",
						uid.String(),
						time.Now(),
						time.Now(),
						fmt.Sprintf("{4|%v|6}", uid),
					)

				mock.ExpectQuery(regexp.QuoteMeta(userOrdersQ)).
					WithArgs(uid, size, (page-1)*size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
			},
		},
		{
			name: "Repo Not Found",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(userOrderCountQ)).
					WithArgs(uid).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.Equal(t, repo2.ErrNotFound, err)
			},
		},
		{
			name: "Repo Internal Error",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(userOrderCountQ)).
					WithArgs(uid).
					WillReturnError(errors.New("internal error"))
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "Empty",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(userOrderCountQ)).
					WithArgs(uid).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectQuery(regexp.QuoteMeta(userOrdersQ)).
					WithArgs(uid, size, (page-1)*size).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"id",
								"status",
								"total_amount",
								"fio",
								"tel",
								"email",
								"address",
								"delivery",
								"payment_method",
								"user_id",
								"created_at",
								"updated_at",
								"order_items",
							},
						),
					)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, int64(0), res.Count)
				assert.Len(t, res.Data, 0)
				assert.Equal(t, 0, res.TotalPages)
				assert.False(t, res.HasNextPage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.ListUserOrders(context.Background(), uid, page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_GetOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	orderID := uint64(12345)

	itemID1 := uuid.New()
	expectedOrder := &model.Order{
		ID:            orderID,
		Status:        "completed",
		TotalAmount:   100.0,
		FIO:           "John Doe",
		Tel:           "1234567890",
		Email:         "john@example.com",
		Address:       "123 Street",
		Delivery:      "delivery",
		PaymentMethod: "credit_card",
		UserID:        uuid.New(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		OrderItems: []*model.OrderItem{
			{
				ID:       1,
				ItemID:   itemID1,
				Quantity: 3,
				Item: model.Item{
					ID:      itemID1,
					Title:   "Item 1",
					Article: "Article 1",
					Price:   10.0,
					Src:     "src1",
					Alt:     "alt1",
					InStock: true,
				},
			},
		},
	}

	tests := []struct {
		name         string
		orderID      uint64
		mockExpect   func()
		expectedResp func(*testing.T, *model.Order, error)
	}{
		{
			name:    "Success",
			orderID: orderID,
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(orderGetQ)).
					WithArgs(orderID).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"id",
								"status",
								"total_amount",
								"fio",
								"tel",
								"email",
								"address",
								"delivery",
								"payment_method",
								"user_id",
								"created_at",
								"updated_at",
								"order_items",
							},
						).
							AddRow(
								expectedOrder.ID,
								expectedOrder.Status,
								expectedOrder.TotalAmount,
								expectedOrder.FIO,
								expectedOrder.Tel,
								expectedOrder.Email,
								expectedOrder.Address,
								expectedOrder.Delivery,
								expectedOrder.PaymentMethod,
								expectedOrder.UserID.String(),
								expectedOrder.CreatedAt,
								expectedOrder.UpdatedAt,
								fmt.Sprintf("{1|%v|3}", itemID1),
							),
					)

				for _, item := range expectedOrder.OrderItems {
					mock.ExpectQuery(
						regexp.QuoteMeta(
							`SELECT 
							 i.id, i.title, i.article, i.price, i.src, i.alt, i.in_stock
							 FROM item i
							 WHERE id = $1`,
						),
					).
						WithArgs(item.ItemID).
						WillReturnRows(
							sqlmock.NewRows(
								[]string{"id", "title", "article", "price", "src", "alt", "in_stock"},
							).
								AddRow(
									item.Item.ID,
									item.Item.Title,
									item.Item.Article,
									item.Item.Price,
									item.Item.Src,
									item.Item.Alt,
									item.Item.InStock,
								),
						)
				}
			},
			expectedResp: func(t *testing.T, res *model.Order, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name:    "Repo Not Found",
			orderID: orderID,
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(orderGetQ)).
					WithArgs(orderID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res *model.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.Equal(t, repo2.ErrNotFound, err)
			},
		},
		{
			name:    "Repo Internal Error",
			orderID: orderID,
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(orderGetQ)).
					WithArgs(orderID).
					WillReturnError(errors.New("internal error"))
			},
			expectedResp: func(t *testing.T, res *model.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.GetOrder(context.Background(), tt.orderID)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_CreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()
	orderID := uint64(12345)

	order := &model.Order{
		FIO:           "John Doe",
		Tel:           "1234567890",
		Email:         "john@example.com",
		Address:       "123 Street",
		Delivery:      "delivery",
		PaymentMethod: "credit_card",
		UserID:        uid,
		OrderItems: []*model.OrderItem{
			{
				ItemID:   uuid.New(),
				Quantity: 3,
				Item: model.Item{
					ID:    uuid.New(),
					Price: 10.0,
				},
			},
			{
				ItemID:   uuid.New(),
				Quantity: 5,
				Item: model.Item{
					ID:    uuid.New(),
					Price: 20.0,
				},
			},
		},
	}

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, uint64, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(orderCreateQ)).
					WithArgs(
						model.OrderStatusPending,
						130.0,
						order.FIO,
						order.Tel,
						order.Email,
						order.Address,
						order.Delivery,
						order.PaymentMethod,
						uid,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(orderID))

				for _, item := range order.OrderItems {
					mock.ExpectExec(regexp.QuoteMeta(orderItemCreateQ)).
						WithArgs(orderID, item.ItemID, item.Quantity).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, res uint64, err error) {
				require.NoError(t, err)
				assert.Equal(t, orderID, res)
			},
		},
		{
			name: "Create Order Error",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(orderCreateQ)).
					WithArgs(
						model.OrderStatusPending,
						130.0, // total amount
						order.FIO,
						order.Tel,
						order.Email,
						order.Address,
						order.Delivery,
						order.PaymentMethod,
						uid,
					).
					WillReturnError(errors.New("create order error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, res uint64, err error) {
				assert.Error(t, err)
				assert.Equal(t, uint64(0), res)
			},
		},
		{
			name: "Create Order Item Error",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(orderCreateQ)).
					WithArgs(
						model.OrderStatusPending,
						130.0, // total amount
						order.FIO,
						order.Tel,
						order.Email,
						order.Address,
						order.Delivery,
						order.PaymentMethod,
						uid,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(orderID))

				mock.ExpectExec(regexp.QuoteMeta(orderItemCreateQ)).
					WithArgs(orderID, order.OrderItems[0].ItemID, order.OrderItems[0].Quantity).
					WillReturnError(errors.New("create order item error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, res uint64, err error) {
				assert.Error(t, err)
				assert.Equal(t, uint64(0), res)
			},
		},
		{
			name: "Commit Error",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(orderCreateQ)).
					WithArgs(
						model.OrderStatusPending,
						130.0, // total amount
						order.FIO,
						order.Tel,
						order.Email,
						order.Address,
						order.Delivery,
						order.PaymentMethod,
						uid,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(orderID))

				for _, item := range order.OrderItems {
					mock.ExpectExec(regexp.QuoteMeta(orderItemCreateQ)).
						WithArgs(orderID, item.ItemID, item.Quantity).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResp: func(t *testing.T, res uint64, err error) {
				assert.Error(t, err)
				assert.Equal(t, uint64(0), res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.CreateOrder(context.Background(), uid, order)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_UpdateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	orderID := uint64(12345)

	order := &model.Order{
		FIO:           "John Doe",
		Status:        model.OrderStatusPending,
		Tel:           "1234567890",
		Email:         "john@example.com",
		Address:       "123 Street",
		Delivery:      "delivery",
		PaymentMethod: "credit_card",
		OrderItems:    []*model.OrderItem{},
	}

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(orderUpdateQ)).
					WithArgs(
						model.OrderStatusPending,
						0.,
						order.FIO,
						order.Tel,
						order.Email,
						order.Address,
						order.Delivery,
						order.PaymentMethod,
						orderID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Update Order Error",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(orderUpdateQ)).
					WithArgs(
						model.OrderStatusPending,
						0.,
						order.FIO,
						order.Tel,
						order.Email,
						order.Address,
						order.Delivery,
						order.PaymentMethod,
						orderID,
					).
					WillReturnError(errors.New("update order error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "update order error", err.Error())
			},
		},
		{
			name: "Commit Error",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(orderUpdateQ)).
					WithArgs(
						model.OrderStatusPending,
						0.,
						order.FIO,
						order.Tel,
						order.Email,
						order.Address,
						order.Delivery,
						order.PaymentMethod,
						orderID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "commit error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				err := repo.UpdateOrder(context.Background(), orderID, order)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_CancelOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	orderID := uint64(12345)

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(orderCancelQ)).
					WithArgs(model.OrderStatusCancelled, orderID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Repo Not Found",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(orderCancelQ)).
					WithArgs(model.OrderStatusCancelled, orderID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
			},
		},
		{
			name: "Repo Internal Error",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(orderCancelQ)).
					WithArgs(model.OrderStatusCancelled, orderID).
					WillReturnError(errors.New("internal error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "internal error", err.Error())
			},
		},
		{
			name: "Commit Error",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(orderCancelQ)).
					WithArgs(model.OrderStatusCancelled, orderID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "commit error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				err := repo.CancelOrder(context.Background(), orderID)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}
