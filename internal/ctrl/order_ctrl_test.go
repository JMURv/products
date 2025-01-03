package ctrl

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestController_ListOrders(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	ctrl := New(rr, nil)

	tests := []struct {
		name         string
		page         int
		size         int
		filters      map[string]any
		sort         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.PaginatedOrderData, error)
	}{
		{
			name:    "Success",
			page:    1,
			size:    10,
			filters: map[string]any{},
			sort:    "created_at DESC",
			mockExpect: func() {
				rr.EXPECT().ListOrders(gomock.Any(), 1, 10, gomock.Any(), "created_at DESC").Return(
					&model.PaginatedOrderData{
						Data:        []*model.Order{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name:    "RepoInternalError",
			page:    1,
			size:    10,
			filters: map[string]any{},
			sort:    "created_at DESC",
			mockExpect: func() {
				rr.EXPECT().ListOrders(gomock.Any(), 1, 10, gomock.Any(), "created_at DESC").Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.ListOrders(context.Background(), tt.page, tt.size, tt.filters, tt.sort)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ListUserOrders(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	uid := uuid.New()
	expRes := &model.PaginatedOrderData{
		Data:        []*model.Order{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		uid          uuid.UUID
		page         int
		size         int
		mockExpect   func()
		expectedResp func(*testing.T, *model.PaginatedOrderData, error)
	}{
		{
			name: "Cache Hit",
			uid:  uid,
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderUserCacheKey, uid, 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "Repo Success",
			uid:  uid,
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderUserCacheKey, uid, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListUserOrders(
					gomock.Any(),
					uid,
					1,
					10,
				).Return(expRes, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(orderUserCacheKey, uid, 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, expRes, res)
			},
		},
		{
			name: "Repo Not Found",
			uid:  uid,
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderUserCacheKey, uid, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListUserOrders(
					gomock.Any(),
					uid,
					1,
					10,
				).Return(nil, repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name: "Repo Internal Error",
			uid:  uid,
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderUserCacheKey, uid, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListUserOrders(
					gomock.Any(),
					uid,
					1,
					10,
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedOrderData, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.ListUserOrders(context.Background(), tt.uid, tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_GetOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	orderID := uint64(12345)
	expRes := &model.Order{
		ID: orderID,
	}

	tests := []struct {
		name         string
		orderID      uint64
		mockExpect   func()
		expectedResp func(*testing.T, *model.Order, error)
	}{
		{
			name:    "Cache Hit",
			orderID: orderID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderCacheKey, orderID),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.Order, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name:    "Repo Success",
			orderID: orderID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderCacheKey, orderID),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetOrder(
					gomock.Any(),
					orderID,
				).Return(expRes, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(orderCacheKey, orderID),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.Order, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, expRes, res)
			},
		},
		{
			name:    "Repo Not Found",
			orderID: orderID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderCacheKey, orderID),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetOrder(
					gomock.Any(),
					orderID,
				).Return(nil, repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.Order, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name:    "Repo Internal Error",
			orderID: orderID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(orderCacheKey, orderID),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetOrder(
					gomock.Any(),
					orderID,
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.Order, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.GetOrder(context.Background(), tt.orderID)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_CreateOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	uid := uuid.New()
	order := &model.Order{}

	tests := []struct {
		name         string
		uid          uuid.UUID
		order        *model.Order
		mockExpect   func()
		expectedResp func(*testing.T, uint64, error)
	}{
		{
			name:  "Repo Success",
			uid:   uid,
			order: order,
			mockExpect: func() {
				rr.EXPECT().CreateOrder(gomock.Any(), uid, order).Return(uint64(12345), nil).Times(1)
				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					invalidateOrderRelatedCachePattern,
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res uint64, err error) {
				require.NoError(t, err)
				assert.Equal(t, uint64(12345), res)
			},
		},
		{
			name:  "Repo Not Found",
			uid:   uid,
			order: order,
			mockExpect: func() {
				rr.EXPECT().CreateOrder(gomock.Any(), uid, order).Return(uint64(0), repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res uint64, err error) {
				require.Error(t, err)
				assert.Equal(t, uint64(0), res)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name:  "Repo Internal Error",
			uid:   uid,
			order: order,
			mockExpect: func() {
				rr.EXPECT().CreateOrder(gomock.Any(), uid, order).Return(
					uint64(0),
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res uint64, err error) {
				require.Error(t, err)
				assert.Equal(t, uint64(0), res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.CreateOrder(context.Background(), tt.uid, tt.order)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_UpdateOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	orderID := uint64(12345)
	ctx := context.Background()
	newData := &model.Order{
		ID: orderID,
	}

	tests := []struct {
		name         string
		orderID      uint64
		newData      *model.Order
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name:    "Success",
			orderID: orderID,
			newData: newData,
			mockExpect: func() {
				rr.EXPECT().UpdateOrder(gomock.Any(), orderID, newData).Return(nil).Times(1)
				cc.EXPECT().Delete(gomock.Any(), fmt.Sprintf(orderCacheKey, orderID)).Return(nil).Times(1)
				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					invalidateOrderRelatedCachePattern,
				).Return(nil).AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:    "Repo Not Found",
			orderID: orderID,
			newData: newData,
			mockExpect: func() {
				rr.EXPECT().UpdateOrder(gomock.Any(), orderID, newData).Return(repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name:    "Repo Internal Error",
			orderID: orderID,
			newData: newData,
			mockExpect: func() {
				rr.EXPECT().UpdateOrder(gomock.Any(), orderID, newData).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:    "Cache Delete Error",
			orderID: orderID,
			newData: newData,
			mockExpect: func() {
				rr.EXPECT().UpdateOrder(gomock.Any(), orderID, newData).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(orderCacheKey, orderID),
				).Return(errors.New("cache delete error")).Times(1)
				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					invalidateOrderRelatedCachePattern,
				).Return(nil).AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				err := ctrl.UpdateOrder(ctx, tt.orderID, tt.newData)

				tt.expectedResp(t, err)
			},
		)
	}
}

func TestController_CancelOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	orderID := uint64(12345)
	ctx := context.Background()

	tests := []struct {
		name         string
		orderID      uint64
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name:    "Success",
			orderID: orderID,
			mockExpect: func() {
				rr.EXPECT().CancelOrder(gomock.Any(), orderID).Return(nil).Times(1)
				cc.EXPECT().Delete(gomock.Any(), fmt.Sprintf(orderCacheKey, orderID)).Return(nil).Times(1)
				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					invalidateOrderRelatedCachePattern,
				).Return(nil).AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:    "Repo Not Found",
			orderID: orderID,
			mockExpect: func() {
				rr.EXPECT().CancelOrder(gomock.Any(), orderID).Return(repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name:    "Repo Internal Error",
			orderID: orderID,
			mockExpect: func() {
				rr.EXPECT().CancelOrder(gomock.Any(), orderID).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:    "Cache Delete Error",
			orderID: orderID,
			mockExpect: func() {
				rr.EXPECT().CancelOrder(gomock.Any(), orderID).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(orderCacheKey, orderID),
				).Return(errors.New("cache delete error")).Times(1)
				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					invalidateOrderRelatedCachePattern,
				).Return(nil).AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				err := ctrl.CancelOrder(ctx, tt.orderID)

				tt.expectedResp(t, err)
			},
		)
	}
}
