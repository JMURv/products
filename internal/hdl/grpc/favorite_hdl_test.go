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

func TestHandler_ListFavorites(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.UuidMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.FavoriteListMsg, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.UuidMsg{Uuid: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.FavoriteListMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid UUID",
			req:        &pb.UuidMsg{Uuid: "invalid-uuid"},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.FavoriteListMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().ListFavorites(gomock.Any(), gomock.Any()).Return(
					[]*model.Favorite{}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FavoriteListMsg, err error) {
				expectedRes := &pb.FavoriteListMsg{
					Data: []*pb.FavoriteMsg{},
				}
				assert.Equal(t, expectedRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().ListFavorites(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FavoriteListMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListFavorites(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_AddToFavorites(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.UserAndItemIds
		mockExpect   func()
		expectedResp func(*testing.T, *pb.FavoriteMsg, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.UserAndItemIds{UserId: "", ItemId: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.FavoriteMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid user UUID",
			req:        &pb.UserAndItemIds{UserId: "invalid-uuid", ItemId: uuid.New().String()},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.FavoriteMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid item UUID",
			req:        &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: "invalid-uuid"},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.FavoriteMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&model.Favorite{}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FavoriteMsg, err error) {
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, ctrl.ErrNotFound,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FavoriteMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Already exists",
			req:  &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, ctrl.ErrAlreadyExists,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FavoriteMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.AlreadyExists, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FavoriteMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.AddToFavorites(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_RemoveFromFavorites(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.UserAndItemIds
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.UserAndItemIds{UserId: "", ItemId: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid user UUID",
			req:        &pb.UserAndItemIds{UserId: "invalid-uuid", ItemId: uuid.New().String()},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid item UUID",
			req:        &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: "invalid-uuid"},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().RemoveFromFavorites(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().RemoveFromFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.UserAndItemIds{UserId: uuid.New().String(), ItemId: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().RemoveFromFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.RemoveFromFavorites(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}
