package ctrl

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestController_PromotionSearch(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		query        string
		page         int
		size         int
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name:  "CacheHit",
			query: "test-query",
			page:  1,
			size:  10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(promotionSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name:  "RepoSuccess",
			query: "test-query",
			page:  1,
			size:  10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(promotionSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().PromotionSearch(
					gomock.Any(),
					"test-query",
					1,
					10,
				).Return(&model.PaginatedPromosData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(promotionSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name:  "RepoInternalError",
			query: "test-query",
			page:  1,
			size:  10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(promotionSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().PromotionSearch(
					gomock.Any(),
					"test-query",
					1,
					10,
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.PromotionSearch(context.Background(), tt.query, tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ListPromotions(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		page         int
		size         int
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "CacheHit",
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(PromotionListCacheKey, 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "RepoSuccess",
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(PromotionListCacheKey, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListPromotions(
					gomock.Any(),
					1,
					10,
				).Return(&model.PaginatedPromosData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(PromotionListCacheKey, 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "RepoInternalError",
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(PromotionListCacheKey, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListPromotions(
					gomock.Any(),
					1,
					10,
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.ListPromotions(context.Background(), tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_GetPromotion(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		slug         string
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "CacheHit",
			slug: "test-promotion",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "RepoNotFound",
			slug: "test-promotion",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetPromotion(
					gomock.Any(),
					"test-promotion",
				).Return(nil, repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "RepoSuccess",
			slug: "test-promotion",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetPromotion(
					gomock.Any(),
					"test-promotion",
				).Return(&model.Promotion{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "RepoInternalError",
			slug: "test-promotion",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetPromotion(
					gomock.Any(),
					"test-promotion",
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.GetPromotion(context.Background(), tt.slug)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_CreatePromotion(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		promotion    *model.Promotion
		mockExpect   func()
		expectedResp func(*testing.T, string, error)
	}{
		{
			name: "RepoSuccess",
			promotion: &model.Promotion{
				Title: "Test Promotion",
			},
			mockExpect: func() {
				rr.EXPECT().CreatePromotion(
					gomock.Any(),
					gomock.Any(),
				).Return("test-promotion-slug", nil).Times(1)

				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					gomock.Any(),
				).AnyTimes()
			},
			expectedResp: func(t *testing.T, res string, err error) {
				require.NoError(t, err)
				assert.Equal(t, "test-promotion-slug", res)
			},
		},
		{
			name: "RepoError",
			promotion: &model.Promotion{
				Title: "Test Promotion",
			},
			mockExpect: func() {
				rr.EXPECT().CreatePromotion(
					gomock.Any(),
					gomock.Any(),
				).Return("", errors.New("repo error")).Times(1)
			},
			expectedResp: func(t *testing.T, res string, err error) {
				require.Error(t, err)
				assert.Empty(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.CreatePromotion(context.Background(), tt.promotion)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_UpdatePromotion(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		slug         string
		promotion    *model.Promotion
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "RepoSuccess",
			slug: "test-promotion",
			promotion: &model.Promotion{
				Title: "Updated Promotion",
			},
			mockExpect: func() {
				rr.EXPECT().UpdatePromotion(
					gomock.Any(),
					"test-promotion",
					gomock.Any(),
				).Return(nil).Times(1)

				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
				).Return(nil).Times(1)

				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					gomock.Any(),
				).AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "RepoNotFound",
			slug: "test-promotion",
			promotion: &model.Promotion{
				Title: "Updated Promotion",
			},
			mockExpect: func() {
				rr.EXPECT().UpdatePromotion(
					gomock.Any(),
					"test-promotion",
					gomock.Any(),
				).Return(repo.ErrNotFound).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name: "RepoInternalError",
			slug: "test-promotion",
			promotion: &model.Promotion{
				Title: "Updated Promotion",
			},
			mockExpect: func() {
				rr.EXPECT().UpdatePromotion(
					gomock.Any(),
					"test-promotion",
					gomock.Any(),
				).Return(errors.New("internal error")).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "CacheDeleteError",
			slug: "test-promotion",
			promotion: &model.Promotion{
				Title: "Updated Promotion",
			},
			mockExpect: func() {
				rr.EXPECT().UpdatePromotion(
					gomock.Any(),
					"test-promotion",
					gomock.Any(),
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
				).Return(errors.New("cache delete error")).Times(1)
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

				err := ctrl.UpdatePromotion(context.Background(), tt.slug, tt.promotion)

				tt.expectedResp(t, err)
			},
		)
	}
}

func TestController_DeletePromotion(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		slug         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "RepoSuccess",
			slug: "test-promotion",
			mockExpect: func() {
				rr.EXPECT().DeletePromotion(
					gomock.Any(),
					"test-promotion",
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
				).Return(nil).Times(1)

				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					gomock.Any(),
				).AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "RepoNotFound",
			slug: "test-promotion",
			mockExpect: func() {
				rr.EXPECT().DeletePromotion(
					gomock.Any(),
					"test-promotion",
				).Return(repo.ErrNotFound).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name: "RepoInternalError",
			slug: "test-promotion",
			mockExpect: func() {
				rr.EXPECT().DeletePromotion(
					gomock.Any(),
					"test-promotion",
				).Return(errors.New("internal error")).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "CacheDeleteError",
			slug: "test-promotion",
			mockExpect: func() {
				rr.EXPECT().DeletePromotion(
					gomock.Any(),
					"test-promotion",
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(promotionCacheKey, "test-promotion"),
				).Return(errors.New("cache delete error")).Times(1)
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

				err := ctrl.DeletePromotion(context.Background(), tt.slug)

				tt.expectedResp(t, err)
			},
		)
	}
}

func TestController_ListPromotionItems(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		slug         string
		page         int
		size         int
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "CacheHit",
			slug: "test-promotion",
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(PromotionItemsCacheKey, "test-promotion", 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "RepoSuccess",
			slug: "test-promotion",
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(PromotionItemsCacheKey, "test-promotion", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListPromotionItems(
					gomock.Any(),
					"test-promotion",
					1,
					10,
				).Return(&model.PaginatedPromoItemsData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(PromotionItemsCacheKey, "test-promotion", 1, 10),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "RepoInternalError",
			slug: "test-promotion",
			page: 1,
			size: 10,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(PromotionItemsCacheKey, "test-promotion", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListPromotionItems(
					gomock.Any(),
					"test-promotion",
					1,
					10,
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.ListPromotionItems(context.Background(), tt.slug, tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}
