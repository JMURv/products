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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	invalidUUID  = "invalid-uuid"
	validUUID    = uuid.New()
	validTitle   = "TestTitle"
	validDesc    = "TestDescription"
	invalidPrice = 0
	validPrice   = float64(299)
	validSrc     = "TestPath"
)
var validItm = &model.Item{
	ID:          validUUID,
	Title:       validTitle,
	Description: validDesc,
	Price:       validPrice,
	Src:         validSrc,
}

func TestHandler_ItemSearch(t *testing.T) {
	const uri = "/api/item/search"
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
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedItemsData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=10&page=invalid",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedItemsData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testquery", 1, 10).Return(
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
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testquery", 1, 10).Return(
					&model.PaginatedItemsData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				h.itemSearch(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListItems(t *testing.T) {
	const uri = "/api/items"
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
			name:    "InvalidPage - fallback to default",
			method:  http.MethodGet,
			url:     uri + "?page=invalid&size=10",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListItems(gomock.Any(), 1, 10).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidSize - fallback to default",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=invalid",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListItems(gomock.Any(), 1, consts.DefaultPageSize).Return(
					&model.PaginatedItemsData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				mctrl.EXPECT().ListItems(gomock.Any(), 1, 10).Return(nil, errors.New("internal error")).Times(1)
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
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListItems(gomock.Any(), 1, 10).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				h.ListItems(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_GetItem(t *testing.T) {
	const uri = "/api/item/"
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
			name:       "InvalidUUID",
			method:     http.MethodGet,
			url:        uri + "invalid-uuid",
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID")
			},
		},
		{
			name:    "ItemNotFound",
			method:  http.MethodGet,
			url:     uri + uuid.Nil.String(),
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().GetItemByUUID(gomock.Any(), uuid.MustParse(uuid.Nil.String())).Return(
					nil,
					repo.ErrNotFound,
				).Times(1)
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
			url:     uri + uuid.Nil.String(),
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().GetItemByUUID(gomock.Any(), uuid.MustParse(uuid.Nil.String())).Return(
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
			url:     uri + uuid.Nil.String(),
			body:    nil,
			resType: &model.Item{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().GetItemByUUID(gomock.Any(), uuid.MustParse(uuid.Nil.String())).Return(
					&model.Item{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.Item)
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
				h.GetItem(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_CreateItem(t *testing.T) {
	const uri = "/api/item"
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
			body:    &model.Item{},
			resType: &utils.ErrorResponse{},
			status:  http.StatusBadRequest,
			mockExpect: func() {
				mctrl.EXPECT().CreateItem(gomock.Any(), gomock.Any()).Times(0)
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
			body:    validItm,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().CreateItem(gomock.Any(), gomock.Any()).Return(
					validUUID,
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
			body:    validItm,
			resType: &utils.Response{},
			status:  http.StatusCreated,
			mockExpect: func() {
				mctrl.EXPECT().CreateItem(gomock.Any(), validItm).Return(
					validUUID,
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				log.Println(res)
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Equal(t, validUUID.String(), resp.Data)
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
				h.CreateItem(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_UpdateItem(t *testing.T) {
	const uri = "/api/item/"
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
			name:       "InvalidUUID",
			method:     http.MethodPut,
			url:        uri + invalidUUID,
			body:       validItm,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID")
			},
		},
		{
			name:       "InvalidBody",
			method:     http.MethodPut,
			url:        uri + validUUID.String(),
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
			url:     uri + validUUID.String(),
			body:    &model.Item{},
			resType: &utils.ErrorResponse{},
			status:  http.StatusBadRequest,
			mockExpect: func() {
				mctrl.EXPECT().UpdateItem(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
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
			url:     uri + validUUID.String(),
			body:    validItm,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().UpdateItem(gomock.Any(), validUUID, validItm).Return(
					repo.ErrNotFound,
				).Times(1)
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
			url:     uri + validUUID.String(),
			body:    validItm,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().UpdateItem(gomock.Any(), validUUID, validItm).Return(
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
			method:  http.MethodPut,
			url:     uri + validUUID.String(),
			body:    validItm,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().UpdateItem(gomock.Any(), validUUID, validItm).Return(
					nil,
				).Times(1)
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
				h.UpdateItem(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_DeleteItem(t *testing.T) {
	const uri = "/api/item/"
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
			name:       "InvalidUUID",
			method:     http.MethodDelete,
			url:        uri + invalidUUID,
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID")
			},
		},
		{
			name:    "NotFound",
			method:  http.MethodDelete,
			url:     uri + validUUID.String(),
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().DeleteItem(gomock.Any(), validUUID).Return(repo.ErrNotFound).Times(1)
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
			url:     uri + validUUID.String(),
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().DeleteItem(gomock.Any(), validUUID).Return(errors.New("internal error")).Times(1)
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
			url:     uri + validUUID.String(),
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusNoContent,
			mockExpect: func() {
				mctrl.EXPECT().DeleteItem(gomock.Any(), validUUID).Return(nil).Times(1)
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
				h.DeleteItem(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListRelatedItems(t *testing.T) {
	const uri = "/api/item/related/"
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
			name:       "InvalidUUID",
			method:     http.MethodGet,
			url:        uri + invalidUUID,
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID")
			},
		},
		{
			name:    "NotFound",
			method:  http.MethodGet,
			url:     uri + validUUID.String(),
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().ListRelatedItems(gomock.Any(), validUUID).Return(nil, repo.ErrNotFound).Times(1)
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
			url:     uri + validUUID.String(),
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListRelatedItems(gomock.Any(), validUUID).Return(
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
			url:     uri + validUUID.String(),
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListRelatedItems(
					gomock.Any(),
					validUUID,
				).Return([]*model.RelatedProduct{{ItemID: validUUID}}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				resp, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Len(t, resp.Data, 1)
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
				h.listRelatedItems(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_HitItems(t *testing.T) {
	const uri = "/api/items/hit"
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
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().HitItems(gomock.Any(), 1, 10).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=invalid",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().HitItems(gomock.Any(), 1, consts.DefaultPageSize).Return(
					&model.PaginatedItemsData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				mctrl.EXPECT().HitItems(gomock.Any(), 1, 10).Return(nil, errors.New("internal error")).Times(1)
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
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().HitItems(gomock.Any(), 1, 10).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				h.HitItems(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_RecItems(t *testing.T) {
	const uri = "/api/items/rec"
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
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().RecItems(gomock.Any(), 1, 10).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=invalid",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().RecItems(gomock.Any(), 1, consts.DefaultPageSize).Return(
					&model.PaginatedItemsData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				mctrl.EXPECT().RecItems(gomock.Any(), 1, 10).Return(nil, errors.New("internal error")).Times(1)
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
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().RecItems(gomock.Any(), 1, 10).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				h.RecItems(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListCategoryItems(t *testing.T) {
	const uri = "/api/category/items/"
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
			url:     uri + "test-category?page=invalid&size=10",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryItems(
					gomock.Any(),
					"test-category",
					1,
					10,
					gomock.Any(),
					gomock.Any(),
				).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "test-category?page=1&size=invalid",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryItems(
					gomock.Any(),
					"test-category",
					1,
					consts.DefaultPageSize,
					gomock.Any(),
					gomock.Any(),
				).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
				require.True(t, ok)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri + "test-category?page=1&size=10",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryItems(
					gomock.Any(),
					"test-category",
					1,
					10,
					gomock.Any(),
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
			name:    "Success",
			method:  http.MethodGet,
			url:     uri + "test-category?page=1&size=10",
			body:    nil,
			resType: &model.PaginatedItemsData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListCategoryItems(
					gomock.Any(),
					"test-category",
					1,
					10,
					gomock.Any(),
					gomock.Any(),
				).Return(&model.PaginatedItemsData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
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
				h.listCategoryItems(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ItemAttrSearch(t *testing.T) {
	const uri = "/api/item/attr/search"
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
			resType: &model.PaginatedItemAttrData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ItemAttrSearch(gomock.Any(), "testquery", 10, 1).Return(
					&model.PaginatedItemAttrData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemAttrData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "?q=testquery&size=10&page=invalid",
			body:    nil,
			resType: &model.PaginatedItemAttrData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ItemAttrSearch(gomock.Any(), "testquery", 10, 1).Return(
					&model.PaginatedItemAttrData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemAttrData)
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
				mctrl.EXPECT().ItemAttrSearch(gomock.Any(), "testquery", 10, 1).Return(
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
			resType: &model.PaginatedItemAttrData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ItemAttrSearch(gomock.Any(), "testquery", 10, 1).Return(
					&model.PaginatedItemAttrData{},
					nil,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemAttrData)
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
				h.itemAttrSearch(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}
