package http

import (
	"errors"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	"github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/pkg/model"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RegisterFavoriteRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc(
		"/api/favorite", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				middlewareFunc(h.listFavorites, h.authMiddleware)
			case http.MethodPost:
				middlewareFunc(h.addToFavorites, h.authMiddleware)
			case http.MethodDelete:
				middlewareFunc(h.removeFromFavorites, h.authMiddleware)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			}
		},
	)
}

func (h *Handler) listFavorites(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "favorites.listFavorites.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	uid, err := uuid.Parse(r.Context().Value("uid").(string))
	if err != nil {
		c = http.StatusUnauthorized
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.ListFavorites(r.Context(), uid)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to find favorites", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to list favorites", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) addToFavorites(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusCreated
	const op = "favorites.addToFavorites.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	uid, err := uuid.Parse(r.Context().Value("uid").(string))
	if err != nil {
		c = http.StatusUnauthorized
		utils.ErrResponse(w, c, err)
		return
	}

	req := &model.Favorite{}
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if req.ItemID == uuid.Nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, validation.ErrMissingUUID)
		return
	}

	res, err := h.ctrl.AddToFavorites(r.Context(), uid, req.ItemID)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		c = http.StatusConflict
		zap.L().Debug("failed to add to favorites", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("failed to find item", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to add to favorites", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) removeFromFavorites(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusNoContent
	const op = "favorites.removeFromFavorites.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	uid, err := uuid.Parse(r.Context().Value("uid").(string))
	if err != nil {
		c = http.StatusUnauthorized
		utils.ErrResponse(w, c, err)
		return
	}

	req := &model.Favorite{}
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if req.ItemID == uuid.Nil {
		c = http.StatusBadRequest
		utils.ErrResponse(w, c, validation.ErrMissingUUID)
		return
	}

	err = h.ctrl.RemoveFromFavorites(r.Context(), uid, req.ItemID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		c = http.StatusNotFound
		zap.L().Debug("Failed to find favorite", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		zap.L().Debug("failed to remove from favorites", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}
