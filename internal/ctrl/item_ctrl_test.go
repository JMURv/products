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

func TestController_ListCategoryItems(t *testing.T) {
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
		filters      map[string]any
		sort         string
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name:    "CacheHit",
			slug:    "test-category",
			page:    1,
			size:    10,
			filters: map[string]any{},
			sort:    "name",
			mockExpect: func() {
				cc.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name:    "RepoNotFound",
			slug:    "non-existent-category",
			page:    1,
			size:    10,
			filters: map[string]any{},
			sort:    "name",
			mockExpect: func() {
				cc.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("cache miss")).
					Times(1)
				rr.EXPECT().ListCategoryItems(
					gomock.Any(),
					"non-existent-category",
					1,
					10,
					map[string]any{},
					"name",
				).Return(
					nil,
					repo.ErrNotFound,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Equal(t, repo.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name:    "RepoSuccess",
			slug:    "test-category",
			page:    1,
			size:    10,
			filters: map[string]any{},
			sort:    "name",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListCategoryItems(
					gomock.Any(),
					"test-category",
					1,
					10,
					map[string]any{},
					"name",
				).Return(&model.PaginatedItemsData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name:    "RepoInternalError",
			slug:    "test-category",
			page:    1,
			size:    10,
			filters: map[string]any{},
			sort:    "name",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListCategoryItems(gomock.Any(), "test-category", 1, 10, map[string]any{}, "name").Return(
					nil,
					errors.New("internal error"),
				).Times(1)
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

				res, err := ctrl.ListCategoryItems(context.Background(), tt.slug, tt.page, tt.size, tt.filters, tt.sort)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ItemAttrSearch(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		query        string
		size         int
		page         int
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name:  "CacheHit",
			query: "test-query",
			size:  10,
			page:  1,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(itemAttrSearchCacheKey, "test-query", 1, 10),
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
			size:  10,
			page:  1,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(itemAttrSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ItemAttrSearch(
					gomock.Any(),
					"test-query",
					10,
					1,
				).Return(&model.PaginatedItemAttrData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(itemAttrSearchCacheKey, "test-query", 1, 10),
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
			size:  10,
			page:  1,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(itemAttrSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ItemAttrSearch(
					gomock.Any(),
					"test-query",
					10,
					1,
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

				res, err := ctrl.ItemAttrSearch(context.Background(), tt.query, tt.size, tt.page)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ItemSearch(t *testing.T) {
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
					fmt.Sprintf(itemSearchCacheKey, "test-query", 1, 10),
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
					fmt.Sprintf(itemSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ItemSearch(
					gomock.Any(),
					"test-query",
					1,
					10,
				).Return(&model.PaginatedItemsData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(itemSearchCacheKey, "test-query", 1, 10),
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
					fmt.Sprintf(itemSearchCacheKey, "test-query", 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ItemSearch(
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

				res, err := ctrl.ItemSearch(context.Background(), tt.query, tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ListRelatedItems(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	uid := uuid.New()
	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		uid          uuid.UUID
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "CacheHit",
			uid:  uid,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(relatedItemCacheKey, uid),
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
			uid:  uid,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(relatedItemCacheKey, uid),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListRelatedItems(
					gomock.Any(),
					uid,
				).Return([]*model.RelatedProduct{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(relatedItemCacheKey, uid),
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
			uid:  uid,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(relatedItemCacheKey, uid),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListRelatedItems(
					gomock.Any(),
					uid,
				).Return(nil, repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Equal(t, repo.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "RepoInternalError",
			uid:  uid,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(relatedItemCacheKey, uid),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListRelatedItems(
					gomock.Any(),
					uid,
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

				res, err := ctrl.ListRelatedItems(context.Background(), tt.uid)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ListItemsByLabel(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	label := "test-label"
	page := 1
	size := 10
	cacheKey := fmt.Sprintf(itemLabelCacheKey, label, page, size)

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "CacheHit",
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					cacheKey,
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
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					cacheKey,
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListItemsByLabel(
					gomock.Any(),
					label,
					page,
					size,
				).Return(&model.PaginatedItemsData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					cacheKey,
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
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					cacheKey,
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListItemsByLabel(
					gomock.Any(),
					label,
					page,
					size,
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

				res, err := ctrl.ListItemsByLabel(context.Background(), label, page, size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_ListItems(t *testing.T) {
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
					fmt.Sprintf(itemListCacheKey, 1, 10),
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
					fmt.Sprintf(itemListCacheKey, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListItems(
					gomock.Any(),
					1,
					10,
				).Return(&model.PaginatedItemsData{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(itemListCacheKey, 1, 10),
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
					fmt.Sprintf(itemListCacheKey, 1, 10),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().ListItems(
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

				res, err := ctrl.ListItems(context.Background(), tt.page, tt.size)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_GetItemByUUID(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	itemID := uuid.New()
	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		uid          uuid.UUID
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "CacheHit",
			uid:  itemID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(itemCacheKey, itemID),
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
			uid:  itemID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(itemCacheKey, itemID),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetItemByUUID(
					gomock.Any(),
					itemID,
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
			uid:  itemID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(itemCacheKey, itemID),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetItemByUUID(
					gomock.Any(),
					itemID,
				).Return(&model.Item{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(itemCacheKey, itemID),
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
			uid:  itemID,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(itemCacheKey, itemID),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetItemByUUID(
					gomock.Any(),
					itemID,
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

				res, err := ctrl.GetItemByUUID(context.Background(), tt.uid)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_CreateItem(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	itemID := uuid.New()
	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		item         *model.Item
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "RepoSuccess",
			item: &model.Item{},
			mockExpect: func() {
				rr.EXPECT().CreateItem(
					gomock.Any(),
					gomock.Any(),
				).Return(itemID, nil).Times(1)

				cc.EXPECT().
					InvalidateKeysByPattern(gomock.Any(), gomock.Any()).
					AnyTimes()
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, res)
			},
		},
		{
			name: "RepoError",
			item: &model.Item{},
			mockExpect: func() {
				rr.EXPECT().CreateItem(
					gomock.Any(),
					gomock.Any(),
				).Return(uuid.Nil, errors.New("repo error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.Error(t, err)
				assert.Equal(t, uuid.Nil, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.CreateItem(context.Background(), tt.item)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_UpdateItem(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		uid          uuid.UUID
		item         *model.Item
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "RepoSuccess",
			uid:  uuid.New(),
			item: &model.Item{Title: "Test Item"},
			mockExpect: func() {
				rr.EXPECT().UpdateItem(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)

				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)

				cc.EXPECT().
					InvalidateKeysByPattern(gomock.Any(), gomock.Any()).
					AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "RepoNotFound",
			uid:  uuid.New(),
			item: &model.Item{Title: "Test Item"},
			mockExpect: func() {
				rr.EXPECT().UpdateItem(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name: "RepoInternalError",
			uid:  uuid.New(),
			item: &model.Item{Title: "Test Item"},
			mockExpect: func() {
				rr.EXPECT().UpdateItem(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "CacheDeleteError",
			uid:  uuid.New(),
			item: &model.Item{Title: "Test Item"},
			mockExpect: func() {
				rr.EXPECT().UpdateItem(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
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
				err := ctrl.UpdateItem(context.Background(), tt.uid, tt.item)
				tt.expectedResp(t, err)
			},
		)
	}
}

func TestController_DeleteItem(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		uid          uuid.UUID
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "RepoSuccess",
			uid:  uuid.New(),
			mockExpect: func() {
				rr.EXPECT().DeleteItem(
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)

				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)

				cc.EXPECT().
					InvalidateKeysByPattern(gomock.Any(), gomock.Any()).
					AnyTimes()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "RepoNotFound",
			uid:  uuid.New(),
			mockExpect: func() {
				rr.EXPECT().DeleteItem(
					gomock.Any(),
					gomock.Any(),
				).Return(repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name: "RepoInternalError",
			uid:  uuid.New(),
			mockExpect: func() {
				rr.EXPECT().DeleteItem(
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "CacheDeleteError",
			uid:  uuid.New(),
			mockExpect: func() {
				rr.EXPECT().DeleteItem(
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)

				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
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

				err := ctrl.DeleteItem(context.Background(), tt.uid)

				tt.expectedResp(t, err)
			},
		)
	}
}
