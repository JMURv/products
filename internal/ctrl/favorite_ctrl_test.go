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

func TestController_ListFavorites(t *testing.T) {
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
					fmt.Sprintf(favoriteCacheKey, uid),
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
					fmt.Sprintf(favoriteCacheKey, uid),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetFavorites(
					gomock.Any(),
					gomock.Any(),
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
			uid:  uid,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(favoriteCacheKey, uid),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetFavorites(
					gomock.Any(),
					gomock.Any(),
				).Return([]*model.Favorite{}, nil).Times(1)
				cc.EXPECT().Set(
					gomock.Any(),
					consts.DefaultCacheTime,
					fmt.Sprintf(favoriteCacheKey, uid),
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
			uid:  uid,
			mockExpect: func() {
				cc.EXPECT().GetToStruct(
					gomock.Any(),
					fmt.Sprintf(favoriteCacheKey, uid),
					gomock.Any(),
				).Return(errors.New("cache miss")).Times(1)
				rr.EXPECT().GetFavorites(
					gomock.Any(),
					gomock.Any(),
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

				res, err := ctrl.ListFavorites(context.Background(), tt.uid)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_AddToFavorites(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	uID := uuid.New()
	iID := uuid.New()
	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		uid          uuid.UUID
		itemID       uuid.UUID
		mockExpect   func()
		expectedResp func(*testing.T, *model.Favorite, error)
	}{
		{
			name:   "RepoSuccess",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().AddToFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(&model.Favorite{}, nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(favoriteCacheKey, uID),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name:   "RepoAlreadyExists",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().AddToFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, repo.ErrAlreadyExists).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrAlreadyExists, err)
				assert.Nil(t, res)
			},
		},
		{
			name:   "RepoNotFound",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().AddToFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, repo.ErrNotFound).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name:   "RepoInternalError",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().AddToFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("internal error")).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				require.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				res, err := ctrl.AddToFavorites(context.Background(), tt.uid, tt.itemID)

				tt.expectedResp(t, res, err)
			},
		)
	}
}

func TestController_RemoveFromFavorites(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	uID := uuid.New()
	iID := uuid.New()
	rr := mocks.NewMockAppRepo(mock)
	cc := mocks.NewMockCacheService(mock)
	ctrl := New(rr, cc)

	tests := []struct {
		name         string
		uid          uuid.UUID
		itemID       uuid.UUID
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name:   "RepoSuccess",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().RemoveFromFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(favoriteCacheKey, uID),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:   "RepoNotFound",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().RemoveFromFavorites(
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
			name:   "RepoInternalError",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().RemoveFromFavorites(
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
			name:   "CacheDeleteError",
			uid:    uID,
			itemID: iID,
			mockExpect: func() {
				rr.EXPECT().RemoveFromFavorites(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
				cc.EXPECT().Delete(
					gomock.Any(),
					fmt.Sprintf(favoriteCacheKey, uID),
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

				err := ctrl.RemoveFromFavorites(context.Background(), tt.uid, tt.itemID)

				tt.expectedResp(t, err)
			},
		)
	}
}
