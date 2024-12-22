package middleware

import (
	"github.com/JMURv/par-pro/products/internal/ctrl"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"go.uber.org/zap"
	"net/http"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					zap.L().Error("panic", zap.Any("err", err))
					utils.ErrResponse(w, http.StatusInternalServerError, ctrl.ErrInternalError)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
