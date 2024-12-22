package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ListFavorites(t *testing.T) {
	const uri = "/api/favorite"
	mock := gomock.NewController(t)
	defer mock.Finish()

	uid := uuid.New()
	invalidCtx := context.WithValue(context.Background(), "uid", uid.String()+"1")
	validCtx := context.WithValue(context.Background(), "uid", uid.String())

	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
		ctx          context.Context
		name         string
		method       string
		url          string
		body         any
		resType      any
		status       int
		mockExpect   func()
		expectedResp func(*testing.T, any)
	}{
		{
			ctx:        invalidCtx,
			name:       "InvalidUID",
			method:     http.MethodGet,
			url:        uri,
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusUnauthorized,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID")
			},
		},
		{
			ctx:     validCtx,
			name:    "NotFound",
			method:  http.MethodGet,
			url:     uri,
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().ListFavorites(
					gomock.Any(),
					uid,
				).Return(nil, repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, repo.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri,
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListFavorites(
					gomock.Any(),
					uid,
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "Success",
			method:  http.MethodGet,
			url:     uri,
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListFavorites(
					gomock.Any(),
					uid,
				).Return([]*model.Favorite{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.NotNil(t, resp.Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				req := httptest.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(tt.ctx)

				w := httptest.NewRecorder()
				h.listFavorites(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_AddToFavorites(t *testing.T) {
	const uri = "/api/favorites"
	mock := gomock.NewController(t)
	defer mock.Finish()

	uid := uuid.New()
	invalidCtx := context.WithValue(context.Background(), "uid", uid.String()+"1")
	validCtx := context.WithValue(context.Background(), "uid", uid.String())

	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
		ctx          context.Context
		name         string
		method       string
		url          string
		body         any
		resType      any
		status       int
		mockExpect   func()
		expectedResp func(*testing.T, any)
	}{
		{
			ctx:        invalidCtx,
			name:       "InvalidUID",
			method:     http.MethodPost,
			url:        uri,
			body:       &model.Favorite{ItemID: uuid.New()},
			resType:    &utils.ErrorResponse{},
			status:     http.StatusUnauthorized,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID")
			},
		},
		{
			ctx:        validCtx,
			name:       "InvalidBody",
			method:     http.MethodPost,
			url:        uri,
			body:       "invalid body",
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid character")
			},
		},
		{
			ctx:        validCtx,
			name:       "MissingUUID",
			method:     http.MethodPost,
			url:        uri,
			body:       &model.Favorite{ItemID: uuid.Nil},
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, validation.ErrMissingUUID.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "AlreadyExists",
			method:  http.MethodPost,
			url:     uri,
			body:    &model.Favorite{ItemID: uuid.New()},
			resType: &utils.ErrorResponse{},
			status:  http.StatusConflict,
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(
					gomock.Any(),
					uid,
					gomock.Any(),
				).Return(nil, repo.ErrAlreadyExists).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, repo.ErrAlreadyExists.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "NotFound",
			method:  http.MethodPost,
			url:     uri,
			body:    &model.Favorite{ItemID: uuid.New()},
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(
					gomock.Any(),
					uid,
					gomock.Any(),
				).Return(nil, repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, repo.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "InternalError",
			method:  http.MethodPost,
			url:     uri,
			body:    &model.Favorite{ItemID: uuid.New()},
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(
					gomock.Any(),
					uid,
					gomock.Any(),
				).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "Success",
			method:  http.MethodPost,
			url:     uri,
			body:    &model.Favorite{ItemID: uuid.New()},
			resType: &utils.Response{},
			status:  http.StatusCreated,
			mockExpect: func() {
				mctrl.EXPECT().AddToFavorites(
					gomock.Any(),
					uid,
					gomock.Any(),
				).Return(&model.Favorite{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.NotNil(t, resp.Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				var body *bytes.Reader
				if tt.body != nil {
					bodyBytes, err := json.Marshal(tt.body)
					require.NoError(t, err)
					body = bytes.NewReader(bodyBytes)
				} else {
					body = bytes.NewReader([]byte{})
				}

				req := httptest.NewRequest(tt.method, tt.url, body)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(tt.ctx)

				w := httptest.NewRecorder()
				h.addToFavorites(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_RemoveFromFavorites(t *testing.T) {
	const uri = "/api/favorites"
	mock := gomock.NewController(t)
	defer mock.Finish()

	uid := uuid.New()
	invalidCtx := context.WithValue(context.Background(), "uid", uid.String()+"1")
	validCtx := context.WithValue(context.Background(), "uid", uid.String())

	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
		ctx          context.Context
		name         string
		method       string
		url          string
		body         any
		resType      any
		status       int
		mockExpect   func()
		expectedResp func(*testing.T, any)
	}{
		{
			ctx:        invalidCtx,
			name:       "InvalidUID",
			method:     http.MethodDelete,
			url:        uri,
			body:       &model.Favorite{ItemID: uuid.New()},
			resType:    &utils.ErrorResponse{},
			status:     http.StatusUnauthorized,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID")
			},
		},
		{
			ctx:        validCtx,
			name:       "InvalidBody",
			method:     http.MethodDelete,
			url:        uri,
			body:       "invalid body",
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid character")
			},
		},
		{
			ctx:        validCtx,
			name:       "MissingUUID",
			method:     http.MethodDelete,
			url:        uri,
			body:       &model.Favorite{ItemID: uuid.Nil},
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, validation.ErrMissingUUID.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "NotFound",
			method:  http.MethodDelete,
			url:     uri,
			body:    &model.Favorite{ItemID: uuid.New()},
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().RemoveFromFavorites(
					gomock.Any(),
					uid,
					gomock.Any(),
				).Return(repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, repo.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "InternalError",
			method:  http.MethodDelete,
			url:     uri,
			body:    &model.Favorite{ItemID: uuid.New()},
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().RemoveFromFavorites(
					gomock.Any(),
					uid,
					gomock.Any(),
				).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			ctx:     validCtx,
			name:    "Success",
			method:  http.MethodDelete,
			url:     uri,
			body:    &model.Favorite{ItemID: uuid.New()},
			resType: &utils.Response{},
			status:  http.StatusNoContent,
			mockExpect: func() {
				mctrl.EXPECT().RemoveFromFavorites(
					gomock.Any(),
					uid,
					gomock.Any(),
				).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Equal(t, "OK", resp.Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				var body *bytes.Reader
				if tt.body != nil {
					bodyBytes, err := json.Marshal(tt.body)
					require.NoError(t, err)
					body = bytes.NewReader(bodyBytes)
				} else {
					body = bytes.NewReader([]byte{})
				}

				req := httptest.NewRequest(tt.method, tt.url, body)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(tt.ctx)

				w := httptest.NewRecorder()
				h.removeFromFavorites(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}
