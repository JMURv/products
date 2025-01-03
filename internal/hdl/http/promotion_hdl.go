package http

import (
	"errors"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	mid "github.com/JMURv/par-pro/products/internal/hdl/http/middleware"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RegisterPromotionRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc(
		"/api/promotions/search", mid.ApplyMiddleware(
			h.promotionSearch, mid.MethodNotAllowed(http.MethodGet),
		),
	)
	mux.HandleFunc(
		"/api/promotions/items/", mid.ApplyMiddleware(
			h.listPromotionItems, mid.MethodNotAllowed(http.MethodGet),
		),
	)

	mux.HandleFunc(
		"/api/promotions", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.listPromotions(w, r)
			case http.MethodPost:
				mid.ApplyMiddleware(h.createPromotion, h.authMiddleware)(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, mid.ErrMethodNotAllowed)
			}
		},
	)

	mux.HandleFunc(
		"/api/promotions/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.getPromotion(w, r)
			case http.MethodPut:
				mid.ApplyMiddleware(h.updatePromotion, h.authMiddleware)(w, r)
			case http.MethodDelete:
				mid.ApplyMiddleware(h.deletePromotion, h.authMiddleware)(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, mid.ErrMethodNotAllowed)
			}
		},
	)
}

func (h *Handler) promotionSearch(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "promo.search.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		utils.SuccessResponse(w, c, []string{})
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 10
	}

	res, err := h.ctrl.PromotionSearch(r.Context(), query, page, size)
	if err != nil {
		zap.L().Debug("failed to search promotions", zap.String("op", op), zap.String("query", query), zap.Error(err))
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) listPromotionItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "promo.listPromotionItems.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = consts.DefaultPageSize
	}

	res, err := h.ctrl.ListPromotionItems(
		r.Context(),
		strings.TrimPrefix(r.URL.Path, "/api/promotions/items/"),
		page,
		size,
	)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list promotion items", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) listPromotions(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "promo.listPromotions.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = consts.DefaultPageSize
	}

	res, err := h.ctrl.ListPromotions(r.Context(), page, size)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list promotions", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) getPromotion(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "promo.getPromotion.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	res, err := h.ctrl.GetPromotion(r.Context(), strings.TrimPrefix(r.URL.Path, "/api/promotions/"))
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to find promotion", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to get promotion", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) createPromotion(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusCreated
	const op = "promo.createPromotion.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	p := &model.Promotion{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if err := validation.ValidatePromotion(p); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.CreatePromotion(r.Context(), p)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to create promotion", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) updatePromotion(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "promo.updatePromotion.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	p := &model.Promotion{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if err := validation.ValidatePromotion(p); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	err := h.ctrl.UpdatePromotion(r.Context(), strings.TrimPrefix(r.URL.Path, "/api/promotions/"), p)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to find promotion", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to update promotion", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}

func (h *Handler) deletePromotion(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusNoContent
	const op = "promo.deletePromotion.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	err := h.ctrl.DeletePromotion(r.Context(), strings.TrimPrefix(r.URL.Path, "/api/promotions/"))
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to find promotion", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to delete promotion", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}
