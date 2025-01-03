package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/JMURv/par-pro/products/pkg/consts"
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

func TestHandler_ListOrders(t *testing.T) {
	const uri = "/api/orders"
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
			name:    "InvalidSize",
			method:  http.MethodGet,
			url:     uri + "?size=invalid",
			body:    nil,
			resType: &model.PaginatedOrderData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListOrders(
					gomock.Any(),
					1,
					consts.DefaultPageSize,
					gomock.Any(),
					gomock.Any(),
				).Return(&model.PaginatedOrderData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedOrderData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidPage",
			method:  http.MethodGet,
			url:     uri + "?size=10&page=invalid",
			body:    nil,
			resType: &model.PaginatedOrderData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListOrders(
					gomock.Any(),
					consts.DefaultPage,
					10,
					gomock.Any(),
					gomock.Any(),
				).Return(&model.PaginatedOrderData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedOrderData)
				require.True(t, ok)
			},
		},
		{
			name:    "InternalError",
			method:  http.MethodGet,
			url:     uri + "?size=10&page=1",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListOrders(gomock.Any(), 1, 10, gomock.Any(), gomock.Any()).Return(
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
			url:     uri + "?size=10&page=1",
			body:    nil,
			resType: &model.PaginatedOrderData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListOrders(
					gomock.Any(),
					1,
					10,
					gomock.Any(),
					gomock.Any(),
				).Return(&model.PaginatedOrderData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedOrderData)
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
				h.listOrders(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_ListUserOrders(t *testing.T) {
	const uri = "/api/orders/user"
	mock := gomock.NewController(t)
	defer mock.Finish()

	validUid := uuid.New()
	invalidUid := uuid.New().String() + "invalid"
	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	h := New(mctrl, nil)

	tests := []struct {
		name         string
		uid          string
		method       string
		url          string
		body         any
		resType      any
		status       int
		mockExpect   func()
		expectedResp func(*testing.T, any)
	}{
		{
			name:       "InvalidToken",
			uid:        invalidUid,
			method:     http.MethodGet,
			url:        uri,
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusUnauthorized,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid UUID length")
			},
		},
		{
			name:    "InvalidSize",
			uid:     validUid.String(),
			method:  http.MethodGet,
			url:     uri + "?size=invalid",
			body:    nil,
			resType: &model.PaginatedOrderData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListUserOrders(
					gomock.Any(),
					gomock.Any(),
					consts.DefaultPage,
					consts.DefaultPageSize,
				).Return(&model.PaginatedOrderData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedOrderData)
				require.True(t, ok)
			},
		},
		{
			name:    "InvalidPage",
			uid:     validUid.String(),
			method:  http.MethodGet,
			url:     uri + "?size=10&page=invalid",
			body:    nil,
			resType: &model.PaginatedOrderData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListUserOrders(
					gomock.Any(),
					gomock.Any(),
					consts.DefaultPage,
					10,
				).Return(&model.PaginatedOrderData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedOrderData)
				require.True(t, ok)
			},
		},
		{
			name:    "InternalError",
			uid:     validUid.String(),
			method:  http.MethodGet,
			url:     uri + "?size=10&page=1",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().ListUserOrders(gomock.Any(), gomock.Any(), 1, 10).Return(
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
			uid:     validUid.String(),
			method:  http.MethodGet,
			url:     uri + "?size=10&page=1",
			body:    nil,
			resType: &model.PaginatedOrderData{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().ListUserOrders(
					gomock.Any(),
					gomock.Any(),
					1,
					10,
				).Return(&model.PaginatedOrderData{}, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				_, ok := res.(*model.PaginatedOrderData)
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

				req = req.WithContext(context.WithValue(ctx, "uid", tt.uid))

				w := httptest.NewRecorder()
				h.listUserOrders(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_GetOrder(t *testing.T) {
	const uri = "/api/order"
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
			name:       "InvalidOrderID",
			method:     http.MethodGet,
			url:        uri + "/invalid",
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid syntax")
			},
		},
		{
			name:    "OrderNotFound",
			method:  http.MethodGet,
			url:     uri + "/12345",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().GetOrder(gomock.Any(), uint64(12345)).Return(nil, ctrl.ErrNotFound).Times(1)
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
			url:     uri + "/12345",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().GetOrder(gomock.Any(), uint64(12345)).Return(nil, errors.New("internal error")).Times(1)
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
			url:     uri + "/12345",
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				order := &model.Order{ID: 12345, Status: "completed"}
				mctrl.EXPECT().GetOrder(gomock.Any(), uint64(12345)).Return(order, nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				res, ok := res.(*utils.Response)
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
				h.getOrder(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_CreateOrder(t *testing.T) {
	const uri = "/api/order"
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mctrl := mocks.NewMockCtrl(mock)
	msso := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, msso)

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
			name:       "InvalidRequestBody",
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
			name:       "ValidationError",
			method:     http.MethodPost,
			url:        uri,
			body:       &model.Order{FIO: "", Tel: "1234567890", Email: "test@example.com", Address: "123 Street"},
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, validation.ErrMissingFIO.Error(), errResp.Error)
			},
		},
		{
			name:   "CreateUserError",
			method: http.MethodPost,
			url:    uri,
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType: &utils.ErrorResponse{},
			status:  http.StatusBadRequest,
			mockExpect: func() {
				msso.EXPECT().ParseClaims(gomock.Any(), gomock.Any()).Return("", nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, ctrl.ErrInternalError.Error())
			},
		},
		{
			name:   "Success",
			method: http.MethodPost,
			url:    uri,
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType: &utils.Response{},
			status:  http.StatusCreated,
			mockExpect: func() {
				msso.EXPECT().ParseClaims(gomock.Any(), gomock.Any()).Return(uuid.NewString(), nil).Times(1)
				mctrl.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(1), nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				response, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.NotNil(t, response)
			},
		},
		{
			name:   "OrderAlreadyExists",
			method: http.MethodPost,
			url:    uri,
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType: &utils.ErrorResponse{},
			status:  http.StatusConflict,
			mockExpect: func() {
				msso.EXPECT().ParseClaims(gomock.Any(), gomock.Any()).Return(uuid.NewString(), nil).Times(1)
				mctrl.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					uint64(0),
					ctrl.ErrAlreadyExists,
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrAlreadyExists.Error(), errResp.Error)
			},
		},
		{
			name:   "InternalError",
			method: http.MethodPost,
			url:    uri,
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				msso.EXPECT().ParseClaims(gomock.Any(), gomock.Any()).Return(uuid.NewString(), nil).Times(1)
				mctrl.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					uint64(0),
					errors.New("internal error"),
				).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrInternalError.Error(), errResp.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				reqBody, err := json.Marshal(tt.body)
				require.NoError(t, err)
				req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+uuid.NewString())
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.createOrder(w, req)

				res := tt.resType
				err = json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_UpdateOrder(t *testing.T) {
	const uri = "/api/order"
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
			name:       "InvalidRequestBody",
			method:     http.MethodPut,
			url:        uri + "/12345",
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
			name:       "ValidationError",
			method:     http.MethodPut,
			url:        uri + "/12345",
			body:       &model.Order{FIO: "", Tel: "1234567890", Email: "test@example.com", Address: "123 Street"},
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, validation.ErrMissingFIO.Error(), errResp.Error)
			},
		},
		{
			name:   "InvalidOrderID",
			method: http.MethodPut,
			url:    uri + "/invalid",
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid syntax")
			},
		},
		{
			name:   "OrderNotFound",
			method: http.MethodPut,
			url:    uri + "/12345",
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().UpdateOrder(gomock.Any(), uint64(12345), gomock.Any()).Return(ctrl.ErrNotFound).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Equal(t, ctrl.ErrNotFound.Error(), errResp.Error)
			},
		},
		{
			name:   "InternalError",
			method: http.MethodPut,
			url:    uri + "/12345",
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().UpdateOrder(
					gomock.Any(),
					uint64(12345),
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
			name:   "Success",
			method: http.MethodPut,
			url:    uri + "/12345",
			body: &model.Order{
				FIO:     "Test User",
				Tel:     "1234567890",
				Email:   "test@example.com",
				Address: "123 Street",
			},
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().UpdateOrder(gomock.Any(), uint64(12345), gomock.Any()).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				response, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Equal(t, "OK", response.Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()

				reqBody, err := json.Marshal(tt.body)
				require.NoError(t, err)
				req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				h.updateOrder(w, req)

				res := tt.resType
				err = json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}

func TestHandler_CancelOrder(t *testing.T) {
	const uri = "/api/order"
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
			name:       "InvalidOrderID",
			method:     http.MethodDelete,
			url:        uri + "/invalid",
			body:       nil,
			resType:    &utils.ErrorResponse{},
			status:     http.StatusBadRequest,
			mockExpect: func() {},
			expectedResp: func(t *testing.T, res any) {
				errResp, ok := res.(*utils.ErrorResponse)
				require.True(t, ok)
				assert.Contains(t, errResp.Error, "invalid syntax")
			},
		},
		{
			name:    "OrderNotFound",
			method:  http.MethodDelete,
			url:     uri + "/12345",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusNotFound,
			mockExpect: func() {
				mctrl.EXPECT().CancelOrder(gomock.Any(), uint64(12345)).Return(ctrl.ErrNotFound).Times(1)
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
			url:     uri + "/12345",
			body:    nil,
			resType: &utils.ErrorResponse{},
			status:  http.StatusInternalServerError,
			mockExpect: func() {
				mctrl.EXPECT().CancelOrder(gomock.Any(), uint64(12345)).Return(errors.New("internal error")).Times(1)
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
			url:     uri + "/12345",
			body:    nil,
			resType: &utils.Response{},
			status:  http.StatusOK,
			mockExpect: func() {
				mctrl.EXPECT().CancelOrder(gomock.Any(), uint64(12345)).Return(nil).Times(1)
			},
			expectedResp: func(t *testing.T, res any) {
				response, ok := res.(*utils.Response)
				require.True(t, ok)
				assert.Equal(t, "OK", response.Data)
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
				h.cancelOrder(w, req)

				res := tt.resType
				err := json.NewDecoder(w.Result().Body).Decode(res)
				assert.Nil(t, err)

				assert.Equal(t, tt.status, w.Result().StatusCode)
				tt.expectedResp(t, res)
			},
		)
	}
}
