package http

import (
	"errors"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RegisterItemRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/item/search", h.itemSearch)
	mux.HandleFunc("/api/item/attr/search", h.itemAttrSearch)
	mux.HandleFunc("/api/item/related/", h.ListRelatedItems)
	mux.HandleFunc("/api/category/items/", h.listCategoryItems)
	mux.HandleFunc("/api/hits", h.HitItems)
	mux.HandleFunc("/api/recs", h.RecItems)

	mux.HandleFunc(
		"/api/item", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.ListItems(w, r)
			case http.MethodPost:
				middlewareFunc(h.CreateItem, h.authMiddleware)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			}
		},
	)

	mux.HandleFunc(
		"/api/item/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.GetItem(w, r)
			case http.MethodPut:
				middlewareFunc(h.UpdateItem, h.authMiddleware)
			case http.MethodDelete:
				middlewareFunc(h.DeleteItem, h.authMiddleware)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			}
		},
	)
}

func (h *Handler) listCategoryItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.listCategoryItems.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	if r.Method != http.MethodGet {
		utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

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
	res, err := h.ctrl.ListCategoryItems(
		r.Context(),
		strings.TrimPrefix(r.URL.Path, "/api/category/items/"),
		page,
		size,
		filters,
		sort,
	)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list category items", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	if r.Method != http.MethodGet {
		utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	if r.Method != http.MethodGet {
		utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

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

	res, err := h.ctrl.ItemSearch(r.Context(), query, page, size)
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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err), zap.String("op", op))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	itemUID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/item/"))
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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	itemUID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/item/"))
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	req := &model.Item{}
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if err = validation.ItemValidation(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	err = h.ctrl.UpdateItem(r.Context(), itemUID, req)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug(
			"failed to found item",
			zap.String("op", op), zap.String("uid", itemUID.String()),
			zap.Error(err),
		)
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug(
			"failed to update item",
			zap.String("op", op), zap.String("uid", itemUID.String()),
			zap.Error(err),
		)
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusNoContent
	const op = "items.DeleteItem.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	itemUID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/item/"))
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
		zap.L().Debug(
			"failed to delete item",
			zap.String("op", op),
			zap.String("uid", itemUID.String()),
			zap.Error(err),
		)
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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	if r.Method != http.MethodGet {
		utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	itemUID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/item/related/"))
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
		zap.L().Debug(
			"failed to list related items",
			zap.String("op", op),
			zap.String("uid", itemUID.String()),
			zap.Error(err),
		)
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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	if r.Method != http.MethodGet {
		utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

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

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("panic", zap.Any("err", err))
			c = http.StatusInternalServerError
			utils.ErrResponse(w, c, ctrl.ErrInternalError)
		}
	}()

	if r.Method != http.MethodGet {
		utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

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
