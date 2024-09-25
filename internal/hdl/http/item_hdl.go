package http

import (
	"errors"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func RegisterItemRoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/api/item/search", h.itemSearch).Methods(http.MethodGet)
	r.HandleFunc("/api/item/attr/search", h.itemAttrSearch).Methods(http.MethodGet)
	r.HandleFunc("/api/item", h.ListItems).Methods(http.MethodGet)
	r.HandleFunc("/api/item", middlewareFunc(h.CreateItem, h.authMiddleware)).Methods(http.MethodPost)
	r.HandleFunc("/api/item/{uid}", h.GetItem).Methods(http.MethodGet)
	r.HandleFunc("/api/item/{uid}", middlewareFunc(h.UpdateItem, h.authMiddleware)).Methods(http.MethodPut)
	r.HandleFunc("/api/item/{uid}", middlewareFunc(h.DeleteItem, h.authMiddleware)).Methods(http.MethodDelete)
	r.HandleFunc("/api/item/{uid}/related", h.ListRelatedItems).Methods(http.MethodGet)
	r.HandleFunc("/api/category/{slug}/items", h.listCategoryItems).Methods(http.MethodGet)

	r.HandleFunc("/api/hits", h.HitItems).Methods(http.MethodGet)
	r.HandleFunc("/api/recs", h.RecItems).Methods(http.MethodGet)
}

func (h *Handler) listCategoryItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.listCategoryItems.handler"
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

	sort := r.URL.Query().Get("sort")

	filters := utils.ParseFiltersByURL(r)
	res, err := h.ctrl.ListCategoryItems(r.Context(), mux.Vars(r)["slug"], page, size, filters, sort)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list category items", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) itemAttrSearch(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.itemAttrSearch.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		utils.SuccessResponse(w, c, []string{})
		return
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	res, err := h.ctrl.ItemAttrSearch(r.Context(), query, size, page)
	if err != nil {
		zap.L().Debug("failed to search attributes", zap.String("op", op), zap.String("query", query), zap.Error(err))
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) itemSearch(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.search.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		utils.SuccessResponse(w, c, []string{})
		return
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	res, err := h.ctrl.ItemSearch(r.Context(), query, size, page)
	if err != nil {
		zap.L().Debug("failed to search items", zap.String("op", op), zap.String("query", query), zap.Error(err))
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.ListItems.handler"
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

	res, err := h.ctrl.ListItems(r.Context(), page, size)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list items", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.GetItem.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	itemUID, err := uuid.Parse(mux.Vars(r)["uid"])
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.GetItemByUUID(r.Context(), itemUID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to get item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to get item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusCreated
	const op = "items.CreateItem.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	req := &model.Item{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if err := validation.ItemValidation(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.CreateItem(r.Context(), req)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to create item", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.UpdateItem.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	req := &model.Item{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if err := validation.ItemValidation(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	itemUID, err := uuid.Parse(mux.Vars(r)["uid"])
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.UpdateItem(r.Context(), itemUID, req)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to found item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to update item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusNoContent
	const op = "items.DeleteItem.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	itemUID, err := uuid.Parse(mux.Vars(r)["uid"])
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	err = h.ctrl.DeleteItem(r.Context(), itemUID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to found item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to delete item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}

func (h *Handler) ListRelatedItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.ListRelatedItems.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	itemUID, err := uuid.Parse(mux.Vars(r)["uid"])
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	resp, err := h.ctrl.ListRelatedItems(r.Context(), itemUID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to found item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list related items", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, resp)
}

func (h *Handler) HitItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.HitItems.handler"
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

	resp, err := h.ctrl.HitItems(r.Context(), page, size)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list items", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, resp)
}

func (h *Handler) RecItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.RecItems.handler"
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

	resp, err := h.ctrl.RecItems(r.Context(), page, size)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list items", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessPaginatedResponse(w, c, resp)
}
