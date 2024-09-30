package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/par-pro/products/api/pb"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	"github.com/JMURv/par-pro/products/internal/validation"
	"github.com/JMURv/par-pro/products/pkg/model/mapper"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (h *Handler) CategorySearch(ctx context.Context, req *pb.SearchReq) (*pb.PaginatedCategoryRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.CategorySearch.handler"
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

	res, err := h.ctrl.CategorySearch(ctx, q, int(page), int(size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedCategoryRes{
		Data:        mapper.ListCategoryToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) CategoryFiltersSearch(ctx context.Context, req *pb.SearchReq) (*pb.PaginatedFilterRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.CategoryFiltersSearch.handler"
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

	res, err := h.ctrl.CategoryFiltersSearch(ctx, q, int(page), int(size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedFilterRes{
		Data:        mapper.ListFiltersToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) ListCategories(ctx context.Context, req *pb.ListReq) (*pb.PaginatedCategoryRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.ListCategories.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	page, size := req.Page, req.Size
	if page <= 0 || size <= 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ListCategories(ctx, int(page), int(size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedCategoryRes{
		Data:        mapper.ListCategoryToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) ListCategoryFilters(ctx context.Context, req *pb.SlugMsg) (*pb.FilterListRes, error) {
	s, c := time.Now(), codes.OK
	const op = "items.ListCategoryFilters.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Slug == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ListCategoryFilters(ctx, req.Slug)
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.FilterListRes{
		Data: mapper.ListFiltersToProto(res),
	}, nil
}

func (h *Handler) GetCategory(ctx context.Context, req *pb.SlugMsg) (*pb.CategoryMsg, error) {
	s, c := time.Now(), codes.OK
	const op = "items.GetCategory.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Slug == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.GetCategoryBySlug(ctx, req.Slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return mapper.CategoryToProto(res), nil
}

func (h *Handler) CreateCategory(ctx context.Context, req *pb.CategoryMsg) (*pb.CategoryMsg, error) {
	s, c := time.Now(), codes.OK
	const op = "items.CreateCategory.handler"
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

	obj := mapper.CategoryFromProto(req)
	if err := validation.CategoryValidation(obj); err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, err.Error())
	}

	res, err := h.ctrl.CreateCategory(ctx, obj)
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return mapper.CategoryToProto(res), nil
}

func (h *Handler) UpdateCategory(ctx context.Context, req *pb.CategoryWithSlug) (*pb.CategoryMsg, error) {
	s, c := time.Now(), codes.OK
	const op = "items.UpdateCategory.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Slug == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	obj := mapper.CategoryFromProto(req.Category)
	if err := validation.CategoryValidation(obj); err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to validate request", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, err.Error())
	}

	res, err := h.ctrl.UpdateCategory(ctx, req.Slug, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return mapper.CategoryToProto(res), nil
}

func (h *Handler) DeleteCategory(ctx context.Context, req *pb.SlugMsg) (*pb.Empty, error) {
	s, c := time.Now(), codes.OK
	const op = "items.DeleteCategory.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Slug == "" {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	err := h.ctrl.DeleteCategory(ctx, req.Slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.Empty{}, nil
}
