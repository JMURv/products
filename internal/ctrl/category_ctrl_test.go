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

func TestController_CategorySearch(t *testing.T) {
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
					fmt.Sprintf(categorySearchCacheKey, "test-query", 1, 10),
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
					fmt.Sprintf(categorySearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().CategorySearch(
					gomock.Any(),
					"test-query",
					1,
					10,
				).Return(&model.PaginatedCategoryData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(categorySearchCacheKey, "test-query", 1, 10),
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
					fmt.Sprintf(categorySearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().CategorySearch(
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

				res, err := ctrl.CategorySearch(context.Background(), tt.query, tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ListCategories(t *testing.T) {
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
					fmt.Sprintf(categoryListCacheKey, 1, 10),
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
					fmt.Sprintf(categoryListCacheKey, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListCategories(
					gomock.Any(),
					1,
					10,
				).Return(&model.PaginatedCategoryData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(categoryListCacheKey, 1, 10),
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
					fmt.Sprintf(categoryListCacheKey, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListCategories(
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

				res, err := ctrl.ListCategories(context.Background(), tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_GetCategoryBySlug(t *testing.T) {
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
			slug: "test-slug",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
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
			slug: "test-slug",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetCategoryBySlug(
					gomock.Any(),
					"test-slug",
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
			slug: "test-slug",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetCategoryBySlug(
					gomock.Any(),
					"test-slug",
				).Return(&model.Category{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(categoryCacheKey, "test-slug"),
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
			slug: "test-slug",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetCategoryBySlug(
					gomock.Any(),
					"test-slug",
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

				res, err := ctrl.GetCategoryBySlug(context.Background(), tt.slug)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_CreateCategory(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		category     *model.Category
		mockExpect   func()
		expectedResp func(*testing.T, string, error)
	}{
		{
			name: "RepoSuccess",
			category: &model.Category{
				Title: "Test Category",
			},
			mockExpect: func() {
				rr.EXPECT().CreateCategory(
					gomock.Any(),
					gomock.Any(),
				).Return("test-category-slug", nil).Times(1)

				cc.EXPECT().InvalidateKeysByPattern(
					gomock.Any(),
					gomock.Any(),
				).AnyTimes()
			},
			expectedResp: func(t *testing.T, res string, err error) {
				require.NoError(t, err)
				assert.Equal(t, "test-category-slug", res)
			},
		},
		{
			name: "RepoError",
			category: &model.Category{
				Title: "Test Category",
			},
			mockExpect: func() {
				rr.EXPECT().CreateCategory(
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

				res, err := ctrl.CreateCategory(context.Background(), tt.category)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_UpdateCategory(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		slug         string
		category     *model.Category
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "RepoSuccess",
			slug: "test-slug",
			category: &model.Category{
				Title: "Updated Category",
			},
			mockExpect: func() {
				rr.EXPECT().UpdateCategory(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)

				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
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
			slug: "test-slug",
			category: &model.Category{
				Title: "Updated Category",
			},
			mockExpect: func() {
				rr.EXPECT().UpdateCategory(
					gomock.Any(),
					gomock.Any(),
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
			slug: "test-slug",
			category: &model.Category{
				Title: "Updated Category",
			},
			mockExpect: func() {
				rr.EXPECT().UpdateCategory(
					gomock.Any(),
					gomock.Any(),
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
			slug: "test-slug",
			category: &model.Category{
				Title: "Updated Category",
			},
			mockExpect: func() {
				rr.EXPECT().UpdateCategory(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
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

				err := ctrl.UpdateCategory(context.Background(), tt.slug, tt.category)

				tt.expectedResp(t, err)
			},
		)
	}
}

func TestController_DeleteCategory(t *testing.T) {
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
			slug: "test-slug",
			mockExpect: func() {
				rr.EXPECT().DeleteCategory(
					gomock.Any(),
					"test-slug",
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
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
			slug: "test-slug",
			mockExpect: func() {
				rr.EXPECT().DeleteCategory(
					gomock.Any(),
					"test-slug",
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
			slug: "test-slug",
			mockExpect: func() {
				rr.EXPECT().DeleteCategory(
					gomock.Any(),
					"test-slug",
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
			slug: "test-slug",
			mockExpect: func() {
				rr.EXPECT().DeleteCategory(
					gomock.Any(),
					"test-slug",
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(categoryCacheKey, "test-slug"),
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

				err := ctrl.DeleteCategory(context.Background(), tt.slug)

				tt.expectedResp(t, err)
			},
		)
	}
}

func TestController_CategoryFiltersSearch(t *testing.T) {
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
					fmt.Sprintf(categoryFiltersSearchCacheKey, "test-query", 1, 10),
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
					fmt.Sprintf(categoryFiltersSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().CategoryFiltersSearch(
					gomock.Any(),
					"test-query",
					1,
					10,
				).Return(&model.PaginatedFilterData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(categoryFiltersSearchCacheKey, "test-query", 1, 10),
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
					fmt.Sprintf(categoryFiltersSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().CategoryFiltersSearch(
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

				res, err := ctrl.CategoryFiltersSearch(context.Background(), tt.query, tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ListCategoryFilters(t *testing.T) {
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
			slug: "test-slug",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(categoryFiltersListCacheKey, "test-slug"),
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
			slug: "test-slug",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(categoryFiltersListCacheKey, "test-slug"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListCategoryFilters(
					gomock.Any(),
					"test-slug",
				).Return([]*model.Filter{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(categoryFiltersListCacheKey, "test-slug"),
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
			slug: "test-slug",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(categoryFiltersListCacheKey, "test-slug"),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListCategoryFilters(
					gomock.Any(),
					"test-slug",
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

				res, err := ctrl.ListCategoryFilters(context.Background(), tt.slug)

				tt.expectedResp(t, res, err)
			},
		)
	}
}
