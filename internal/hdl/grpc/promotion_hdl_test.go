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

func TestHandler_PromotionSearch(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedPromoRes{
		Data:        []*pb.PromoMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.SearchReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedPromoRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SearchReq{Query: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().PromotionSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedPromosData{
						Data:        []*model.Promotion{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SearchReq{Query: "testquery", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().PromotionSearch(gomock.Any(), "testquery", 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.PromotionSearch(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListPromotions(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedPromoRes{
		Data:        []*pb.PromoMsg{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedPromoRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ListReq{Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListPromotions(gomock.Any(), 1, 10).Return(
					&model.PaginatedPromosData{
						Data:        []*model.Promotion{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ListReq{Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListPromotions(gomock.Any(), 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListPromotions(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_ListPromotionItems(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	expRes := &pb.PaginatedPromoItemsRes{
		Data:        []*pb.PromoItem{},
		Count:       10,
		TotalPages:  2,
		CurrentPage: 1,
		HasNextPage: true,
	}

	tests := []struct {
		name         string
		req          *pb.ListPromotionItemsReq
		mockExpect   func()
		expectedResp func(*testing.T, *pb.PaginatedPromoItemsRes, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.ListPromotionItemsReq{Slug: "", Page: 0, Size: 0},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoItemsRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.ListPromotionItemsReq{Slug: "test-promo", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListPromotionItems(gomock.Any(), "test-promo", 1, 10).Return(
					&model.PaginatedPromoItemsData{
						Data:        []*model.PromotionItem{},
						Count:       10,
						TotalPages:  2,
						CurrentPage: 1,
						HasNextPage: true,
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoItemsRes, err error) {
				assert.Equal(t, expRes, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.ListPromotionItemsReq{Slug: "test-promo", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListPromotionItems(gomock.Any(), "test-promo", 1, 10).Return(
					nil, ctrl.ErrNotFound,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoItemsRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.ListPromotionItemsReq{Slug: "test-promo", Page: 1, Size: 10},
			mockExpect: func() {
				mctrl.EXPECT().ListPromotionItems(gomock.Any(), "test-promo", 1, 10).Return(
					nil, errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PaginatedPromoItemsRes, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.ListPromotionItems(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_GetPromotion(t *testing.T) {
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
		expectedResp func(*testing.T, *pb.PromoMsg, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.SlugMsg{Slug: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.PromoMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req:  &pb.SlugMsg{Slug: "test-promo"},
			mockExpect: func() {
				mctrl.EXPECT().GetPromotion(gomock.Any(), "test-promo").Return(
					&model.Promotion{
						Slug:  "test-promo",
						Title: "Test Promotion",
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PromoMsg, err error) {
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.SlugMsg{Slug: "test-promo"},
			mockExpect: func() {
				mctrl.EXPECT().GetPromotion(gomock.Any(), "test-promo").Return(nil, ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PromoMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SlugMsg{Slug: "test-promo"},
			mockExpect: func() {
				mctrl.EXPECT().GetPromotion(gomock.Any(), "test-promo").Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.PromoMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := h.GetPromotion(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_CreatePromotion(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()
	promoSlug := "test-promo"

	tests := []struct {
		name         string
		req          *pb.PromoMsg
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
			name:       "Validation error",
			req:        &pb.PromoMsg{Slug: "", Title: ""},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.SlugMsg, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req: &pb.PromoMsg{
				Slug:        "test-promo",
				Title:       "Test Promotion",
				Description: "Test Description",
				Src:         "test-src",
			},
			mockExpect: func() {
				mctrl.EXPECT().CreatePromotion(gomock.Any(), gomock.Any()).Return(promoSlug, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.SlugMsg, err error) {
				assert.Equal(t, &pb.SlugMsg{Slug: promoSlug}, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req: &pb.PromoMsg{
				Slug:        "test-promo",
				Title:       "Test Promotion",
				Description: "Test Description",
				Src:         "test-src",
			},
			mockExpect: func() {
				mctrl.EXPECT().CreatePromotion(gomock.Any(), gomock.Any()).Return(
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
				res, err := h.CreatePromotion(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_UpdatePromotion(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

	ctx := context.Background()

	tests := []struct {
		name         string
		req          *pb.PromoWithSlug
		mockExpect   func()
		expectedResp func(*testing.T, *pb.Empty, error)
	}{
		{
			name:       "Invalid request",
			req:        &pb.PromoWithSlug{Slug: "", Data: nil},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name:       "Validation error",
			req:        &pb.PromoWithSlug{Slug: "test-promo", Data: &pb.PromoMsg{Title: ""}},
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Success",
			req: &pb.PromoWithSlug{
				Slug: "test-promo",
				Data: &pb.PromoMsg{
					Slug:        "test-promo",
					Title:       "Test Promotion",
					Description: "Test Description",
					Src:         "test-src",
				},
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdatePromotion(gomock.Any(), "test-promo", gomock.Any()).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req: &pb.PromoWithSlug{
				Slug: "test-promo",
				Data: &pb.PromoMsg{
					Slug:        "test-promo",
					Title:       "Test Promotion",
					Description: "Test Description",
					Src:         "test-src",
				},
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdatePromotion(
					gomock.Any(),
					"test-promo",
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
			req: &pb.PromoWithSlug{
				Slug: "test-promo",
				Data: &pb.PromoMsg{
					Slug:        "test-promo",
					Title:       "Test Promotion",
					Description: "Test Description",
					Src:         "test-src",
				},
			},
			mockExpect: func() {
				mctrl.EXPECT().UpdatePromotion(
					gomock.Any(),
					"test-promo",
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
				res, err := h.UpdatePromotion(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestHandler_DeletePromotion(t *testing.T) {
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
			req:  &pb.SlugMsg{Slug: "test-promo"},
			mockExpect: func() {
				mctrl.EXPECT().DeletePromotion(gomock.Any(), "test-promo").Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, codes.OK, status.Code(err))
			},
		},
		{
			name: "Not found",
			req:  &pb.SlugMsg{Slug: "test-promo"},
			mockExpect: func() {
				mctrl.EXPECT().DeletePromotion(gomock.Any(), "test-promo").Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res *pb.Empty, err error) {
				assert.Nil(t, res)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Internal error",
			req:  &pb.SlugMsg{Slug: "test-promo"},
			mockExpect: func() {
				mctrl.EXPECT().DeletePromotion(gomock.Any(), "test-promo").Return(errors.New("internal error")).Times(1)
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
				res, err := h.DeletePromotion(ctx, tt.req)
				tt.expectedResp(t, res, err)
			},
		)
	}
}
