package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var validPromotion = &model.Promotion{
	Title:       "Test Promotion",
	Description: "Test Description",
	Src:         "test-src",
}

func TestHandler_PromotionSearch(t *testing.T) {
	const uri = "/api/promotion/search"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
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
			name:       "QueryTooShort",
			method:     http.MethodGet,
			url:        uri + "?q=ab",
			body:       nil,
			resType:    &utils.Response{},
			status:     http.StatusOK,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Empty(t, resp.Data)
			},
		},
		{
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=invalid",
			body:    nil,
			resType: &model.PaginatedPromosData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().PromotionSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedPromosData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromosData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=10&page=invalid",
			body:    nil,
			resType: &model.PaginatedPromosData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().PromotionSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedPromosData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromosData)
				require.True(t, ok)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=10&page=1",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().PromotionSearch(gomock.Any(), "testquery", 1, 10).Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			name:    "Success",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=10&page=1",
			body:    nil,
			resType: &model.PaginatedPromosData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().PromotionSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedPromosData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromosData)
				require.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				req := httptest.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.promotionSearch(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListPromotionItems(t *testing.T) {
	const uri = "/api/promotions/items/"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
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
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "promotion-slug?page=invalid&size=10",
			body:    nil,
			resType: &model.PaginatedPromoItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotionItems(
					gomock.Any(),
					"promotion-slug",
					1,
					10,
				).Return(&model.PaginatedPromoItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromoItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "promotion-slug?page=1&size=invalid",
			body:    nil,
			resType: &model.PaginatedPromoItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotionItems(
					gomock.Any(),
					"promotion-slug",
					1,
					consts.DefaultPageSize,
				).Return(&model.PaginatedPromoItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromoItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri + "promotion-slug?page=1&size=10",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotionItems(gomock.Any(), "promotion-slug", 1, 10).Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			name:    "Success",
			method:  http.MethodGet,
			url:     uri + "promotion-slug?page=1&size=10",
			body:    nil,
			resType: &model.PaginatedPromoItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotionItems(
					gomock.Any(),
					"promotion-slug",
					1,
					10,
				).Return(&model.PaginatedPromoItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromoItemsData)
				require.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				req := httptest.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.listPromotionItems(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListPromotions(t *testing.T) {
	const uri = "/api/promotions"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
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
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "?page=invalid&size=10",
			body:    nil,
			resType: &model.PaginatedPromosData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotions(gomock.Any(), 1, 10).Return(&model.PaginatedPromosData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromosData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=invalid",
			body:    nil,
			resType: &model.PaginatedPromosData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotions(
					gomock.Any(),
					1,
					consts.DefaultPageSize,
				).Return(&model.PaginatedPromosData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromosData)
				require.True(t, ok)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=10",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotions(gomock.Any(), 1, 10).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			name:    "Success",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=10",
			body:    nil,
			resType: &model.PaginatedPromosData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListPromotions(gomock.Any(), 1, 10).Return(&model.PaginatedPromosData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedPromosData)
				require.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				req := httptest.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.listPromotions(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_GetPromotion(t *testing.T) {
	const uri = "/api/promotions/"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
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
			name:    "NotFound",
			method:  http.MethodGet,
			url:     uri + "invalid-slug",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().GetPromotion(gomock.Any(), "invalid-slug").Return(nil, ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri + "test-slug",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().GetPromotion(gomock.Any(), "test-slug").Return(
					nil,
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			name:    "Success",
			method:  http.MethodGet,
			url:     uri + "test-slug",
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().GetPromotion(gomock.Any(), "test-slug").Return(
					&model.Promotion{
						Slug:  "123",
						Title: "Test Promotion",
					}, nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*utils.Response)
				require.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				req := httptest.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.getPromotion(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_CreatePromotion(t *testing.T) {
	const uri = "/api/promotions"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
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
			name:    "ValidationError",
			method:  http.MethodPost,
			url:     uri,
			body:    &model.Promotion{},
			resType: &utils.ErrorResponse{},
			status:  http.StatusBadRequest,
			mockExpect: func() {
				mctrl.EXPECT().CreatePromotion(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "missing")
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodPost,
			url:     uri,
			body:    validPromotion,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().CreatePromotion(gomock.Any(), gomock.Any()).Return(
					"",
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			name:    "Success",
			method:  http.MethodPost,
			url:     uri,
			body:    validPromotion,
			resType: &utils.Response{},
			status:  http.StatusCreated,
			mockExpect: func() {
				mctrl.EXPECT().CreatePromotion(gomock.Any(), gomock.Any()).Return(validPromotion.Slug, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Equal(t, validPromotion.Slug, resp.Data)
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
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.createPromotion(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_UpdatePromotion(t *testing.T) {
	const uri = "/api/promotions/"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
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
			name:       "InvalidBody",
			method:     http.MethodPut,
			url:        uri + "test-promotion",
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
			name:    "ValidationError",
			method:  http.MethodPut,
			url:     uri + "test-promotion",
			body:    &model.Promotion{},
			resType: &utils.ErrorResponse{},
			status:  http.StatusBadRequest,
			mockExpect: func() {
				mctrl.EXPECT().UpdatePromotion(gomock.Any(), "test-promotion", gomock.Any()).Times(0)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "missing")
			},
		},
		{
			name:    "NotFound",
			method:  http.MethodPut,
			url:     uri + "non-existent-promotion",
			body:    validPromotion,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().UpdatePromotion(
					gomock.Any(),
					"non-existent-promotion",
					validPromotion,
				).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodPut,
			url:     uri + "test-promotion",
			body:    validPromotion,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().UpdatePromotion(
					gomock.Any(),
					"test-promotion",
					validPromotion,
				).Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			name:    "Success",
			method:  http.MethodPut,
			url:     uri + "test-promotion",
			body:    validPromotion,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().UpdatePromotion(gomock.Any(), "test-promotion", validPromotion).Return(nil).Times(1)
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
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.updatePromotion(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_DeletePromotion(t *testing.T) {
	const uri = "/api/promotions/"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
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
			name:    "NotFound",
			method:  http.MethodDelete,
			url:     uri + "invalid-slug",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().DeletePromotion(gomock.Any(), "invalid-slug").Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodDelete,
			url:     uri + "test-slug",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().DeletePromotion(gomock.Any(), "test-slug").Return(errors.New("internal error")).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
		{
			name:    "Success",
			method:  http.MethodDelete,
			url:     uri + "test-slug",
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusNoContent,
			mockExpect: func() {
				mctrl.EXPECT().DeletePromotion(gomock.Any(), "test-slug").Return(nil).Times(1)
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

				req := httptest.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.deletePromotion(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}
