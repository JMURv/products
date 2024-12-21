package http

import (
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

func TestHandler_SearchItem(t *testing.T) {
	const uri = "/api/item/search"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	ssoCtrl := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, ssoCtrl)

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
			name:       "Q < 3",
			method:     http.MethodGet,
			url:        uri + "?page=0&size=0&q=",
			body:       nil,
			resType:    &model.PaginatedItemsData{},
			status:     http.StatusOK,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedItemsData)
				require.True(t, ok)
			},
		},
		{
			name:       "MethodNotAllowed",
			method:     http.MethodPost,
			url:        uri + "?page=0&size=0&q=",
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusMethodNotAllowed,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
			},
		},
		{
			name:    "Internal Error",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=40&q=testq",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().
					ItemSearch(gomock.Any(), "testq", 1, consts.DefaultPageSize).
					Return(nil, errors.New("internal error")).
					Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, "internal error", errResp.Error)
			},
		},
		{
			name:    "Panic",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=40&q=testq",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testq", 1, 40).Do(
					func(ctx context.Context, page, size int, q string) {
						panic("internal error")
					},
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
			url:     uri + "?page=1&size=40&q=testq",
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ItemSearch(gomock.Any(), "testq", 1, 40).Return(
					&model.PaginatedItemsData{},
					nil,
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
			name:    "Panic",
			method:  http.MethodGet,
			url:     uri + "?page=1&size=10",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListItems(gomock.Any(), 1, 10).Do(
					func(ctx context.Context, page, size int) {
						panic("internal error")
					},
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
