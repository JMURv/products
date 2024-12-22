package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/internal/repo"
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

func TestHandler_CategoryFiltersSearch(t *testing.T) {
	const uri = "/api/category/filters/search"
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
			resType: &model.PaginatedFilterData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().CategoryFiltersSearch(
					gomock.Any(),
					"testquery",
					1,
					10,
				).Return(&model.PaginatedFilterData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedFilterData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=10&page=invalid",
			body:    nil,
			resType: &model.PaginatedFilterData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().CategoryFiltersSearch(
					gomock.Any(),
					"testquery",
					1,
					10,
				).Return(&model.PaginatedFilterData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedFilterData)
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
				mctrl.EXPECT().CategoryFiltersSearch(gomock.Any(), "testquery", 1, 10).Return(
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
			resType: &model.PaginatedFilterData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().CategoryFiltersSearch(
					gomock.Any(),
					"testquery",
					1,
					10,
				).Return(&model.PaginatedFilterData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedFilterData)
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
				h.categoryFiltersSearch(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_CategorySearch(t *testing.T) {
	const uri = "/api/category/search"
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
			resType: &model.PaginatedCategoryData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().CategorySearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedCategoryData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedCategoryData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=10&page=invalid",
			body:    nil,
			resType: &model.PaginatedCategoryData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().CategorySearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedCategoryData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedCategoryData)
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
				mctrl.EXPECT().CategorySearch(gomock.Any(), "testquery", 1, 10).Return(
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
			resType: &model.PaginatedCategoryData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().CategorySearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedCategoryData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedCategoryData)
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
				h.categorySearch(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListCategoryFilters(t *testing.T) {
	const uri = "/api/category/filters/test-category"
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
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri,
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryFilters(gomock.Any(), "test-category").Return(
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
			url:     uri,
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryFilters(gomock.Any(), "test-category").Return(
					[]*model.Filter{},
					nil,
				).Times(1)
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
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.listCategoryFilters(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListCategories(t *testing.T) {
	const uri = "/api/categories"
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
			resType: &model.PaginatedCategoryData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListCategories(gomock.Any(), 1, 10).Return(&model.PaginatedCategoryData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedCategoryData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=invalid",
			body:    nil,
			resType: &model.PaginatedCategoryData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListCategories(
					gomock.Any(),
					1,
					consts.DefaultPageSize,
				).Return(&model.PaginatedCategoryData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedCategoryData)
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
				mctrl.EXPECT().ListCategories(gomock.Any(), 1, 10).Return(nil, errors.New("internal error")).Times(1)
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
			resType: &model.PaginatedCategoryData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListCategories(gomock.Any(), 1, 10).Return(&model.PaginatedCategoryData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedCategoryData)
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
				h.listCategories(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_CreateCategory(t *testing.T) {
	const uri = "/api/category"
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
			body:    &model.Category{},
			resType: &utils.ErrorResponse{},
			status:  http.StatusBadRequest,
			mockExpect: func() {
				mctrl.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Times(0)
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
			body:    &model.Category{Title: "Test Category"},
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Return(
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
			body:    &model.Category{Title: "Test Category"},
			resType: &utils.Response{},
			status:  http.StatusCreated,
			mockExpect: func() {
				mctrl.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Return(
					"slug", nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Equal(t, "slug", resp.Data)
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
				h.createCategory(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_GetCategory(t *testing.T) {
	const uri = "/api/category/"
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
				mctrl.EXPECT().GetCategoryBySlug(gomock.Any(), "invalid-slug").Return(nil, repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, repo.ErrNotFound.Error(), errResp.Error)
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
				mctrl.EXPECT().GetCategoryBySlug(gomock.Any(), "test-slug").Return(
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
				mctrl.EXPECT().GetCategoryBySlug(
					gomock.Any(),
					"test-slug",
				).Return(&model.Category{Title: "Test Category"}, nil).Times(1)
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
				h.getCategory(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_UpdateCategory(t *testing.T) {
	const uri = "/api/category/"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	validCategory := &model.Category{
		Slug:  "updated-category",
		Title: "Updated Category",
	}

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
			url:        uri + "test-category",
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
			url:     uri + "test-category",
			body:    &model.Category{},
			resType: &utils.ErrorResponse{},
			status:  http.StatusBadRequest,
			mockExpect: func() {
				mctrl.EXPECT().UpdateCategory(gomock.Any(), "test-category", gomock.Any()).Times(0)
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
			url:     uri + "non-existent-category",
			body:    validCategory,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().UpdateCategory(
					gomock.Any(),
					"non-existent-category",
					validCategory,
				).Return(repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, repo.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodPut,
			url:     uri + "test-category",
			body:    validCategory,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().UpdateCategory(
					gomock.Any(),
					"test-category",
					validCategory,
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
			url:     uri + "test-category",
			body:    validCategory,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().UpdateCategory(gomock.Any(), "test-category", validCategory).Return(nil).Times(1)
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
				h.updateCategory(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_DeleteCategory(t *testing.T) {
	const uri = "/api/category/"
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
				mctrl.EXPECT().DeleteCategory(gomock.Any(), "invalid-slug").Return(repo.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, repo.ErrNotFound.Error(), errResp.Error)
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
				mctrl.EXPECT().DeleteCategory(gomock.Any(), "test-slug").Return(errors.New("internal error")).Times(1)
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
				mctrl.EXPECT().DeleteCategory(gomock.Any(), "test-slug").Return(nil).Times(1)
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
				h.deleteCategory(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}
