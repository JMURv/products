package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/par-pro/products/api/pb"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestHandler_CategorySearch(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedCategoryRes{
		Data:        []*pb.CategoryMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.SearchReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedCategoryRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SearchReq{Query: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedCategoryRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().CategorySearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedCategoryData{
						Data:        []*model.Category{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedCategoryRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().CategorySearch(gomock.Any(), "testquery", 1, 10).Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedCategoryRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.CategorySearch(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListCategories(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedCategoryRes{
		Data:        []*pb.CategoryMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedCategoryRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ListReq{Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedCategoryRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListCategories(gomock.Any(), 1, 10).Return(
					&model.PaginatedCategoryData{
						Data:        []*model.Category{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedCategoryRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListCategories(gomock.Any(), 1, 10).Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedCategoryRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListCategories(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_GetCategory(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	returnedCategory := &model.Category{
		Slug:            "test-slug",
		Title:           "Test Category",
		ProductQuantity: int64(1),
		Src:             "test-src",
		Alt:             "test-alt",
		ParentSlug:      "test-parent-slug",
	}

	tests := []struct {
		name         string
		req          *pb.SlugMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.CategoryMsg, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SlugMsg{Slug: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.CategoryMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().GetCategoryBySlug(gomock.Any(), "test-slug").Return(
					returnedCategory, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.CategoryMsg, err error) {
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().GetCategoryBySlug(gomock.Any(), "test-slug").Return(
					nil, ctrl.ErrNotFound,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.CategoryMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().GetCategoryBySlug(gomock.Any(), "test-slug").Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.CategoryMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.GetCategory(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_CreateCategory(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.SlugMsg{Slug: "test-slug"}

	tests := []struct {
		name         string
		req          *pb.CategoryMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.SlugMsg, error)
	}{
		{
			name:       "Invalid request",
			req:        nil,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.SlugMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Validation error",
			req: &pb.CategoryMsg{
				Slug:      "",
				Title:     "",
				Src:       "",
				Alt:       "",
				CreatedAt: nil,
				UpdatedAt: nil,
			},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.SlugMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.CategoryMsg{Slug: "test-slug", Title: "Test Category"},
			mockExpect: func() {
				mctrl.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Return("test-slug", nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.SlugMsg, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.CategoryMsg{Slug: "test-slug", Title: "Test Category"},
			mockExpect: func() {
				mctrl.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Return(
					"",
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.SlugMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.CreateCategory(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_UpdateCategory(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.CategoryWithSlug
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.CategoryWithSlug{Slug: "", Category: nil},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Validation error",
			req:        &pb.CategoryWithSlug{Slug: "test-slug", Category: &pb.CategoryMsg{Title: ""}},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.CategoryWithSlug{Slug: "test-slug", Category: &pb.CategoryMsg{Title: "Test Category"}},
			mockExpect: func() {
				mctrl.EXPECT().UpdateCategory(gomock.Any(), "test-slug", gomock.Any()).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.CategoryWithSlug{Slug: "test-slug", Category: &pb.CategoryMsg{Title: "Test Category"}},
			mockExpect: func() {
				mctrl.EXPECT().UpdateCategory(gomock.Any(), "test-slug", gomock.Any()).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.CategoryWithSlug{Slug: "test-slug", Category: &pb.CategoryMsg{Title: "Test Category"}},
			mockExpect: func() {
				mctrl.EXPECT().UpdateCategory(
					gomock.Any(),
					"test-slug",
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
				res, err := h.UpdateCategory(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_DeleteCategory(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.SlugMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SlugMsg{Slug: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().DeleteCategory(gomock.Any(), "test-slug").Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().DeleteCategory(gomock.Any(), "test-slug").Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().DeleteCategory(gomock.Any(), "test-slug").Return(errors.New("internal error")).Times(1)
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
				res, err := h.DeleteCategory(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_CategoryFiltersSearch(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedFilterRes{
		Data:        []*pb.Filter{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.SearchReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedFilterRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SearchReq{Query: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedFilterRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().CategoryFiltersSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedFilterData{
						Data:        []*model.Filter{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedFilterRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().CategoryFiltersSearch(gomock.Any(), "testquery", 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedFilterRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.CategoryFiltersSearch(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListCategoryFilters(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.SlugMsg
		mockExpect   func()
		expectedResp func(*testing.T, *pb.FilterListRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SlugMsg{Slug: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.FilterListRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryFilters(gomock.Any(), "test-slug").Return(
					[]*model.Filter{}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FilterListRes, err error) {
				expectedRes := &pb.FilterListRes{
					Data: []*pb.Filter{},
				}
				assert.Equal(t, expectedRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SlugMsg{Slug: "test-slug"},
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryFilters(gomock.Any(), "test-slug").Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.FilterListRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListCategoryFilters(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}
