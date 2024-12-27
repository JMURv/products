package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodNotAllowed(t *testing.T) {
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)
	handler := MethodNotAllowed(http.MethodGet, http.MethodPost)(testHandler)

	tests := []struct {
		method         string
		expectedStatus int
	}{
		{method: http.MethodGet, expectedStatus: http.StatusOK},
		{method: http.MethodPost, expectedStatus: http.StatusOK},
		{method: http.MethodPut, expectedStatus: http.StatusMethodNotAllowed},
		{method: http.MethodDelete, expectedStatus: http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(
			tt.method, func(t *testing.T) {
				req, err := http.NewRequest(tt.method, "/", nil)
				assert.NoError(t, err)

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				assert.Equal(t, tt.expectedStatus, rr.Code)
			},
		)
	}
}

func TestRecoverPanic(t *testing.T) {
	panicHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		},
	)
	handler := RecoverPanic(panicHandler)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
