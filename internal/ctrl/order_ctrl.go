package ctrl

import (
	"context"
	"errors"
	"fmt"
	repo "github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/consts"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const orderCacheKey = "order:%v"
const orderUserCacheKey = "orders-user:%v:%v:%v"
const invalidateOrderRelatedCachePattern = "orders-*"

type orderRepo interface {
	ListOrders(ctx context.Context, page, size int, filters map[string]any, sort string) (*model.PaginatedOrderData, error)
	ListUserOrders(ctx context.Context, uid uuid.UUID, page, size int) (*model.PaginatedOrderData, error)
	GetOrder(ctx context.Context, orderID uint64) (*model.Order, error)
	CreateOrder(ctx context.Context, uid uuid.UUID, req *model.Order) (uint64, error)
	UpdateOrder(ctx context.Context, orderID uint64, newData *model.Order) error
	CancelOrder(ctx context.Context, orderID uint64) error
}

func (c *Controller) ListOrders(ctx context.Context, page, size int, filters map[string]any, sort string) (*model.PaginatedOrderData, error) {
	const op = "orders.ListOrders.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.ListOrders(ctx, page, size, filters, sort)
	if err != nil {
		zap.L().Debug("Error list orders", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	return res, nil
}

func (c *Controller) ListUserOrders(ctx context.Context, uid uuid.UUID, page, size int) (*model.PaginatedOrderData, error) {
	const op = "orders.ListUserOrders.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.PaginatedOrderData{}
	cacheKey := fmt.Sprintf(orderUserCacheKey, uid, page, size)
	if err := c.cache.GetToStruct(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.ListUserOrders(ctx, uid, page, size)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("Error list user orders", zap.Error(err), zap.String("op", op))
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("Error list user orders", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}
	return res, nil
}

func (c *Controller) GetOrder(ctx context.Context, orderID uint64) (*model.Order, error) {
	const op = "orders.GetOrder.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.Order{}
	cacheKey := fmt.Sprintf(orderCacheKey, orderID)
	if err := c.cache.GetToStruct(ctx, cacheKey, cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetOrder(ctx, orderID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug("Error get order", zap.Error(err), zap.String("op", op))
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug("failed to set to cache", zap.Error(err), zap.String("op", op))
		}
	}
	return res, nil
}

func (c *Controller) CreateOrder(ctx context.Context, uid uuid.UUID, o *model.Order) (uint64, error) {
	const op = "orders.CreateOrder.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.CreateOrder(ctx, uid, o)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("Error create order", zap.Error(err), zap.String("op", op))
		return 0, ErrNotFound
	} else if err != nil {
		zap.L().Debug("Error create order", zap.Error(err), zap.String("op", op))
		return 0, err
	}

	go c.cache.InvalidateKeysByPattern(ctx, invalidateOrderRelatedCachePattern)
	return res, nil
}

func (c *Controller) UpdateOrder(ctx context.Context, orderID uint64, newData *model.Order) error {
	const op = "orders.UpdateOrder.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.UpdateOrder(ctx, orderID, newData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("Error update order", zap.Error(err), zap.String("op", op))
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(orderCacheKey, orderID)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	go c.cache.InvalidateKeysByPattern(ctx, invalidateOrderRelatedCachePattern)
	return nil
}

func (c *Controller) CancelOrder(ctx context.Context, orderID uint64) error {
	const op = "orders.CancelOrder.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.CancelOrder(ctx, orderID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug("Error cancel order", zap.Error(err), zap.String("op", op))
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(orderCacheKey, orderID)); err != nil {
		zap.L().Debug("failed to delete from cache", zap.Error(err), zap.String("op", op))
	}

	go c.cache.InvalidateKeysByPattern(ctx, invalidateOrderRelatedCachePattern)
	return nil
}
