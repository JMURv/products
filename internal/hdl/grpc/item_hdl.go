package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/par-pro/products/api/pb"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/pkg/model/mapper"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (h *Handler) ItemSearch(ctx context.Context, req *pb.SearchReq) (*pb.PaginatedItemRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.ItemSearch.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	q, page, size := req.Query, req.Page, req.Size
	if q == "" || page == 0 || size == 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ItemSearch(ctx, q, int(page), int(size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        mapper.ListItemToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) ItemAttrSearch(ctx context.Context, req *pb.SearchReq) (*pb.PaginatedItemAttrsRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.ItemAttrSearch.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	q, page, size := req.Query, req.Page, req.Size
	if q == "" || page == 0 || size == 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ItemAttrSearch(ctx, q, int(page), int(size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemAttrsRes{
		Data:        mapper.ListItemAttributesToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) ListItems(ctx context.Context, req *pb.ListReq) (*pb.PaginatedItemRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.ListItems.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	page, size := req.Page, req.Size
	if page == 0 || size == 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ListItems(ctx, int(page), int(size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        mapper.ListItemToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) GetItem(ctx context.Context, req *pb.UuidMsg) (*pb.ItemMsg, error) {
	s, c := time.Now(), codes.OK
	const op = "items.GetItem.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Uuid == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	uid, err := uuid.Parse(req.Uuid)
	if err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrParseUUID.Error())
	}

	res, err := h.ctrl.GetItemByUUID(ctx, uid)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return mapper.ItemToProto(res), nil
}

func (h *Handler) CreateItem(ctx context.Context, req *pb.ItemMsg) (*pb.ItemMsg, error) {
	s, c := time.Now(), codes.OK
	const op = "items.CreateItem.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	item := mapper.ItemFromProto(req)
	if err := validation.ItemValidation(item); err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, err.Error())
	}

	res, err := h.ctrl.CreateItem(ctx, item)
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return mapper.ItemToProto(res), nil
}

func (h *Handler) UpdateItem(ctx context.Context, req *pb.ItemWithUid) (*pb.ItemMsg, error) {
	s, c := time.Now(), codes.OK
	const op = "items.UpdateItem.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Uid == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrParseUUID.Error())
	}

	item := mapper.ItemFromProto(req.Item)
	if err := validation.ItemValidation(item); err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to validate request", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, err.Error())
	}

	res, err := h.ctrl.UpdateItem(ctx, uid, item)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return mapper.ItemToProto(res), nil
}

func (h *Handler) DeleteItem(ctx context.Context, req *pb.UuidMsg) (*pb.Empty, error) {
	s, c := time.Now(), codes.OK
	const op = "items.DeleteItem.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Uuid == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	uid, err := uuid.Parse(req.Uuid)
	if err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrParseUUID.Error())
	}

	err = h.ctrl.DeleteItem(ctx, uid)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.Empty{}, nil
}

func (h *Handler) ListRelatedItems(ctx context.Context, req *pb.UuidMsg) (*pb.RelatedItemsList, error) {
	s, c := time.Now(), codes.OK
	const op = "items.ListRelatedItems.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Uuid == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	uid, err := uuid.Parse(req.Uuid)
	if err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrParseUUID.Error())
	}

	res, err := h.ctrl.ListRelatedItems(ctx, uid)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.RelatedItemsList{
		Items: mapper.ListRelatedProductsToProto(res),
	}, nil
}

func (h *Handler) ListCategoryItems(ctx context.Context, req *pb.ListCategoryItemsReq) (*pb.PaginatedItemRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.ListCategoryItems.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.CategorySlug == "" || req.Page <= 0 || req.Size <= 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	filters := make(map[string]any)
	res, err := h.ctrl.ListCategoryItems(ctx, req.CategorySlug, int(req.Page), int(req.Size), filters, req.Sort)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        mapper.ListItemToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) HitItems(ctx context.Context, req *pb.ListReq) (*pb.PaginatedItemRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.HitItems.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Page <= 0 || req.Size <= 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.HitItems(ctx, int(req.Page), int(req.Size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        mapper.ListItemToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) RecItems(ctx context.Context, req *pb.ListReq) (*pb.PaginatedItemRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.RecItems.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Page <= 0 || req.Size <= 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.RecItems(ctx, int(req.Page), int(req.Size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        mapper.ListItemToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}
