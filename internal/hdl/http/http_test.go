package http

import (
	"errors"
	"github.com/JMURv/par-pro/products/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockCtrl(mock)
	ssoctrl := mocks.NewMockSSOSvc(mock)
	h := New(mctrl, ssoctrl)

	tests := []struct {
		name           string
		authHeader     string
		mockGetIDToken func()
		expectedStatus int
	}{
		{
			name:           "MissingAuthorizationHeader",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "InvalidTokenFormat",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "ValidToken",
			authHeader:     "Bearer valid-token",
			expectedStatus: http.StatusOK,
			mockGetIDToken: func() {
				ssoctrl.EXPECT().GetIDByToken(gomock.Any(), "valid-token").Return("user-id", nil).Times(1)
			},
		},
		{
			name:           "InvalidToken",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			mockGetIDToken: func() {
				ssoctrl.EXPECT().GetIDByToken(gomock.Any(), "invalid-token").Return(
					"",
					errors.New("invalid token"),
				).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if tt.mockGetIDToken != nil {
					tt.mockGetIDToken()
				}

				req, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", tt.authHeader)
				rr := httptest.NewRecorder()

				testHandler := http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
					},
				)

				handler := h.authMiddleware(testHandler)
				handler.ServeHTTP(rr, req)

				assert.Equal(t, tt.expectedStatus, rr.Code)
			},
		)
	}
}
