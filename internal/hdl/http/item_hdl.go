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
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RegisterItemRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc(
		"/api/item/search", func(w http.ResponseWriter, r *http.Request) {
			mid.ApplyMiddleware(h.itemSearch, mid.MethodNotAllowed(http.MethodGet))(w, r)
		},
	)
	mux.HandleFunc(
		"/api/item/attr/search", func(w http.ResponseWriter, r *http.Request) {
			mid.ApplyMiddleware(h.itemAttrSearch, mid.MethodNotAllowed(http.MethodGet))(w, r)
		},
	)
	mux.HandleFunc(
		"/api/item/related/", func(w http.ResponseWriter, r *http.Request) {
			mid.ApplyMiddleware(h.listRelatedItems, mid.MethodNotAllowed(http.MethodGet))(w, r)
		},
	)
	mux.HandleFunc(
		"/api/category/items/", func(w http.ResponseWriter, r *http.Request) {
			mid.ApplyMiddleware(h.listCategoryItems, mid.MethodNotAllowed(http.MethodGet))(w, r)
		},
	)
	mux.HandleFunc(
		"/api/item/label/", func(w http.ResponseWriter, r *http.Request) {
			mid.ApplyMiddleware(h.ListItemsByLabel, mid.MethodNotAllowed(http.MethodGet))(w, r)
		},
	)

	mux.HandleFunc(
		"/api/item", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.ListItems(w, r)
			case http.MethodPost:
				mid.ApplyMiddleware(h.CreateItem, h.authMiddleware)(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, mid.ErrMethodNotAllowed)
			}
		},
	)

	mux.HandleFunc(
		"/api/item/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.GetItem(w, r)
			case http.MethodPut:
				mid.ApplyMiddleware(h.UpdateItem, h.authMiddleware)(w, r)
			case http.MethodDelete:
				mid.ApplyMiddleware(h.DeleteItem, h.authMiddleware)(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, mid.ErrMethodNotAllowed)
			}
		},
	)
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

	res, err := h.ctrl.ItemSearch(r.Context(), query, page, size)
	if err != nil {
		zap.L().Debug("failed to search items", zap.String("op", op), zap.String("query", query), zap.Error(err))
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
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
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
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

	itemUID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/item/"))
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.GetItemByUUID(r.Context(), itemUID)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to get item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to get item", zap.String("op", op), zap.String("uid", itemUID.String()), zap.Error(err))
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
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
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
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
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
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
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
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

	itemUID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/item/"))
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	err = h.ctrl.DeleteItem(r.Context(), itemUID)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
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
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}

func (h *Handler) listCategoryItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.listCategoryItems.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = consts.DefaultPage
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

	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		utils.SuccessResponse(w, c, []string{})
		return
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = consts.DefaultPageSize
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = consts.DefaultPage
	}

	res, err := h.ctrl.ItemAttrSearch(r.Context(), query, size, page)
	if err != nil {
		zap.L().Debug("failed to search attributes", zap.String("op", op), zap.String("query", query), zap.Error(err))
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
		return
	}

	utils.SuccessPaginatedResponse(w, c, res)
}

func (h *Handler) listRelatedItems(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.listRelatedItems.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	itemUID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/item/related/"))
	if err != nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, err)
		return
	}

	resp, err := h.ctrl.ListRelatedItems(r.Context(), itemUID)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
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
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
		return
	}

	utils.SuccessResponse(w, c, resp)
}

func (h *Handler) ListItemsByLabel(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "items.HitItems.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	label := r.URL.Query().Get("label")
	if len(label) < 3 {
		utils.SuccessResponse(w, c, []string{})
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = consts.DefaultPage
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = consts.DefaultPageSize
	}

	resp, err := h.ctrl.ListItemsByLabel(r.Context(), label, page, size)
	if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list items", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, ctrl.ErrInternalError)
		return
	}

	utils.SuccessPaginatedResponse(w, c, resp)
}
