package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/ctrl/sso"
	"github.com/JMURv/par-pro/products/internal/hdl"
	mid "github.com/JMURv/par-pro/products/internal/hdl/http/middleware"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	srv  *http.Server
	ctrl hdl.Ctrl
	sso  sso.SSOSvc
}

func New(ctrl hdl.Ctrl, sso sso.SSOSvc) *Handler {
	return &Handler{
		ctrl: ctrl,
		sso:  sso,
	}
}

func (h *Handler) Start(port int) {
	mux := http.NewServeMux()
	RegisterItemRoutes(mux, h)
	RegisterCategoryRoutes(mux, h)
	RegisterPromotionRoutes(mux, h)
	RegisterFavoriteRoutes(mux, h)
	RegisterOrderRoutes(mux, h)
	mux.HandleFunc(
		"/health-check", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, "OK")
		},
	)

	handler := mid.RecoverPanic(mux)
	handler = mid.TracingMiddleware(mux)
	h.srv = &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err := h.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zap.L().Debug("Server error", zap.Error(err))
	}
}

func (h *Handler) Close() error {
	if err := h.srv.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrResponse(w, http.StatusUnauthorized, errors.New("authorization header is missing"))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				utils.ErrResponse(w, http.StatusUnauthorized, errors.New("invalid token format"))
				return
			}

			token, err := h.sso.ParseClaims(r.Context(), tokenStr)
			if err != nil {
				utils.ErrResponse(w, http.StatusUnauthorized, err)
				return
			}
			ctx := context.WithValue(r.Context(), "uid", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
