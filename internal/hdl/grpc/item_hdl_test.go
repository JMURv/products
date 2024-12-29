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

func TestHandler_ItemSearch(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedItemRes{
		Data:        []*pb.ItemMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.SearchReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedItemRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SearchReq{Query: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedItemsData{
						Data:        []*model.Item{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testquery", 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ItemSearch(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ItemAttrSearch(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedItemAttrsRes{
		Data:        []*pb.ItemAttribute{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.SearchReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedItemAttrsRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SearchReq{Query: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemAttrsRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ItemAttrSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedItemAttrData{
						Data:        []*model.ItemAttribute{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemAttrsRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ItemAttrSearch(gomock.Any(), "testquery", 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemAttrsRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ItemAttrSearch(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListItems(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedItemRes{
		Data:        []*pb.ItemMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedItemRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ListReq{Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListItems(gomock.Any(), 1, 10).Return(
					&model.PaginatedItemsData{
						Data:        []*model.Item{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListItems(gomock.Any(), 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListItems(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_GetItem(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.ItemMsg{
		Id:    uuid.New().String(),
		Title: "Test Item",
	}

	tests := []struct {
		name         string
		req          *pb.UuidMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.ItemMsg, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.UuidMsg{Uuid: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.ItemMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid UUID",
			req:        &pb.UuidMsg{Uuid: "invalid-uuid"},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.ItemMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.UuidMsg{Uuid: expRes.Id},
			mockExpect: func() {
				mctrl.EXPECT().GetItemByUUID(gomock.Any(), gomock.Any()).Return(
					&model.Item{
						ID:    uuid.MustParse(expRes.Id),
						Title: "Test Item",
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.ItemMsg, err error) {
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.UuidMsg{Uuid: expRes.Id},
			mockExpect: func() {
				mctrl.EXPECT().GetItemByUUID(gomock.Any(), gomock.Any()).Return(nil, ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.ItemMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.UuidMsg{Uuid: expRes.Id},
			mockExpect: func() {
				mctrl.EXPECT().GetItemByUUID(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.ItemMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.GetItem(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_CreateItem(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	itemUUID := uuid.New()

	tests := []struct {
		name         string
		req          *pb.ItemMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.UuidMsg, error)
	}{
		{
			name:       "Invalid request",
			req:        nil,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.UuidMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Validation error",
			req:        &pb.ItemMsg{Title: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.UuidMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ItemMsg{Title: "Test Item", Description: "Test Description", Price: 100, Src: "test-src"},
			mockExpect: func() {
				mctrl.EXPECT().CreateItem(gomock.Any(), gomock.Any()).Return(itemUUID, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.UuidMsg, err error) {
				assert.Equal(t, &pb.UuidMsg{Uuid: itemUUID.String()}, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ItemMsg{Title: "Test Item", Description: "Test Description", Price: 100, Src: "test-src"},
			mockExpect: func() {
				mctrl.EXPECT().CreateItem(gomock.Any(), gomock.Any()).Return(
					uuid.Nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.UuidMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.CreateItem(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_UpdateItem(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.ItemWithUid
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ItemWithUid{Uid: "", Item: nil},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid UUID",
			req:        &pb.ItemWithUid{Uid: "invalid-uuid", Item: &pb.ItemMsg{Title: "Test Item"}},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Validation error",
			req:        &pb.ItemWithUid{Uid: uuid.New().String(), Item: &pb.ItemMsg{Title: ""}},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req: &pb.ItemWithUid{
				Uid:  uuid.New().String(),
				Item: &pb.ItemMsg{Title: "Test Item", Description: "Test Description", Price: 100, Src: "test-src"},
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdateItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req: &pb.ItemWithUid{
				Uid:  uuid.New().String(),
				Item: &pb.ItemMsg{Title: "Test Item", Description: "Test Description", Price: 100, Src: "test-src"},
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdateItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req: &pb.ItemWithUid{
				Uid:  uuid.New().String(),
				Item: &pb.ItemMsg{Title: "Test Item", Description: "Test Description", Price: 100, Src: "test-src"},
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdateItem(
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
				res, err := h.UpdateItem(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_DeleteItem(t *testing.T) {
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
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.UuidMsg{Uuid: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid UUID",
			req:        &pb.UuidMsg{Uuid: "invalid-uuid"},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().DeleteItem(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().DeleteItem(gomock.Any(), gomock.Any()).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().DeleteItem(gomock.Any(), gomock.Any()).Return(errors.New("internal error")).Times(1)
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
				res, err := h.DeleteItem(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListRelatedItems(t *testing.T) {
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
		expectedResp func(*testing.T, *pb.RelatedItemsList, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.UuidMsg{Uuid: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.RelatedItemsList, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Invalid UUID",
			req:        &pb.UuidMsg{Uuid: "invalid-uuid"},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.RelatedItemsList, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().ListRelatedItems(gomock.Any(), gomock.Any()).Return(
					[]*model.RelatedProduct{}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.RelatedItemsList, err error) {
				expectedRes := &pb.RelatedItemsList{
					Items: []*pb.RelatedProduct{},
				}
				assert.Equal(t, expectedRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().ListRelatedItems(gomock.Any(), gomock.Any()).Return(nil, ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.RelatedItemsList, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.UuidMsg{Uuid: uuid.New().String()},
			mockExpect: func() {
				mctrl.EXPECT().ListRelatedItems(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.RelatedItemsList, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListRelatedItems(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListCategoryItems(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedItemRes{
		Data:        []*pb.ItemMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListCategoryItemsReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedItemRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ListCategoryItemsReq{CategorySlug: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListCategoryItemsReq{CategorySlug: "test-category", Page: 1, Size: 10, Sort: "asc"},
			mockExpect: func() {
				filters := make(map[string]any)
				mctrl.EXPECT().ListCategoryItems(gomock.Any(), "test-category", 1, 10, filters, "asc").Return(
					&model.PaginatedItemsData{
						Data:        []*model.Item{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.ListCategoryItemsReq{CategorySlug: "test-category", Page: 1, Size: 10, Sort: "asc"},
			mockExpect: func() {
				filters := make(map[string]any)
				mctrl.EXPECT().ListCategoryItems(gomock.Any(), "test-category", 1, 10, filters, "asc").Return(
					nil,
					ctrl.ErrNotFound,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ListCategoryItemsReq{CategorySlug: "test-category", Page: 1, Size: 10, Sort: "asc"},
			mockExpect: func() {
				filters := make(map[string]any)
				mctrl.EXPECT().ListCategoryItems(gomock.Any(), "test-category", 1, 10, filters, "asc").Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListCategoryItems(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListItemsByLabel(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedItemRes{
		Data:        []*pb.ItemMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListItemsByLabelReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedItemRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ListItemsByLabelReq{Label: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListItemsByLabelReq{Label: "test-label", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListItemsByLabel(gomock.Any(), "test-label", 1, 10).Return(
					&model.PaginatedItemsData{
						Data:        []*model.Item{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ListItemsByLabelReq{Label: "test-label", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListItemsByLabel(gomock.Any(), "test-label", 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedItemRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListItemsByLabel(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}
