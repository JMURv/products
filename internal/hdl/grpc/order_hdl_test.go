package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/par-pro/products/api/pb"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestHandler_ListOrders(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	ctx := context.Background()
	expRes := &pb.PaginatedOrderRes{
		Data:        []*pb.OrderMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedOrderRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ListReq{Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedOrderRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListOrders(gomock.Any(), 1, 10, gomock.Any(), "").Return(
					&model.PaginatedOrderData{
						Data:        []*model.Order{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedOrderRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListOrders(gomock.Any(), 1, 10, gomock.Any(), "").Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedOrderRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListOrders(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListUserOrders(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	expRes := &pb.PaginatedOrderRes{
		Data:        []*pb.OrderMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListReq
		ctx          context.Context
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedOrderRes, error)
	}{
		{
			name:       "Unauthenticated",
			req:        &pb.ListReq{Page: 1, Size: 10},
			ctx:        context.Background(),
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedOrderRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Unauthenticated, status.Code(err))
			},
		},
		{
			name:       "Invalid Request",
			req:        &pb.ListReq{Page: 0, Size: 0},
			ctx:        context.WithValue(context.Background(), "uid", "valid-uid"),
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedOrderRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListReq{Page: 1, Size: 10},
			ctx:  context.WithValue(context.Background(), "uid", uuid.New().String()),
			mockExpect: func() {
				mctrl.EXPECT().ListUserOrders(gomock.Any(), gomock.Any(), 1, 10).Return(
					&model.PaginatedOrderData{
						Data:        []*model.Order{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedOrderRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal Error",
			req:  &pb.ListReq{Page: 1, Size: 10},
			ctx:  context.WithValue(context.Background(), "uid", uuid.New().String()),
			mockExpect: func() {
				mctrl.EXPECT().ListUserOrders(gomock.Any(), gomock.Any(), 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedOrderRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListUserOrders(tt.ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_GetOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	ctx := context.Background()
	tests := []struct {
		name         string
		req          *pb.Uint64Msg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.OrderMsg, error)
	}{
		{
			name:       "Invalid Request",
			req:        &pb.Uint64Msg{Value: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.OrderMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Order Not Found",
			req:  &pb.Uint64Msg{Value: 12345},
			mockExpect: func() {
				mctrl.EXPECT().GetOrder(gomock.Any(), uint64(12345)).Return(nil, ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.OrderMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal Error",
			req:  &pb.Uint64Msg{Value: 12345},
			mockExpect: func() {
				mctrl.EXPECT().GetOrder(gomock.Any(), uint64(12345)).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.OrderMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.Uint64Msg{Value: 12345},
			mockExpect: func() {
				mctrl.EXPECT().GetOrder(gomock.Any(), uint64(12345)).Return(&model.Order{ID: 12345}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.OrderMsg, err error) {
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.GetOrder(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_CreateOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	tests := []struct {
		name         string
		req          *pb.OrderMsg
		ctx          context.Context
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Uint64Msg, error)
	}{
		{
			name:       "Invalid Request",
			req:        nil,
			ctx:        context.Background(),
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Uint64Msg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Validation Error",
			req:        &pb.OrderMsg{Fio: "", Email: "test@example.com", UserId: uuid.NewString()},
			ctx:        context.WithValue(context.Background(), "uid", uuid.NewString()),
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Uint64Msg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req: &pb.OrderMsg{
				Address: "some-address", Fio: "Test User", Tel: "some-tel", Email: "test@example.com",
				UserId: uuid.NewString(),
			},
			ctx: context.WithValue(context.Background(), "uid", uuid.NewString()),
			mockExpect: func() {
				mctrl.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(12345), nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Uint64Msg, err error) {
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal Error",
			req: &pb.OrderMsg{
				Address: "some-address", Fio: "Test User", Tel: "some-tel", Email: "test@example.com",
				UserId: uuid.NewString(),
			},
			ctx: context.WithValue(context.Background(), "uid", uuid.NewString()),
			mockExpect: func() {
				mctrl.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					uint64(0),
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Uint64Msg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.CreateOrder(tt.ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_UpdateOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.OrderMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid Request",
			req:        nil,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Validation Error",
			req:        &pb.OrderMsg{Id: 12345, Fio: "", Email: "test@example.com", UserId: uuid.NewString()},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Order Not Found",
			req: &pb.OrderMsg{
				Id:      uint64(12345),
				Address: "some-address", Fio: "Test User", Tel: "some-tel", Email: "test@example.com",
				UserId: uuid.NewString(),
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdateOrder(gomock.Any(), uint64(12345), gomock.Any()).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal Error",
			req: &pb.OrderMsg{
				Id:      uint64(12345),
				Address: "some-address", Fio: "Test User", Tel: "some-tel", Email: "test@example.com",
				UserId: uuid.NewString(),
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdateOrder(
					gomock.Any(),
					uint64(12345),
					gomock.Any(),
				).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
		{
			name: "Success",
			req: &pb.OrderMsg{
				Id:      uint64(12345),
				Address: "some-address", Fio: "Test User", Tel: "some-tel", Email: "test@example.com",
				UserId: uuid.NewString(),
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdateOrder(gomock.Any(), uint64(12345), gomock.Any()).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.UpdateOrder(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_CancelOrder(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.Uint64Msg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid Request",
			req:        nil,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Order Not Found",
			req:  &pb.Uint64Msg{Value: 12345},
			mockExpect: func() {
				mctrl.EXPECT().CancelOrder(gomock.Any(), uint64(12345)).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal Error",
			req:  &pb.Uint64Msg{Value: 12345},
			mockExpect: func() {
				mctrl.EXPECT().CancelOrder(gomock.Any(), uint64(12345)).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.Uint64Msg{Value: 12345},
			mockExpect: func() {
				mctrl.EXPECT().CancelOrder(gomock.Any(), uint64(12345)).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.CancelOrder(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}
