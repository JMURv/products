package grpc

import (
	"context"
	pb "github.com/JMURv/par-pro/products/api/pb"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	"github.com/JMURv/par-pro/products/internal/validation"
	utils "github.com/JMURv/par-pro/products/pkg/utils/grpc"
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
		zap.L().Debug("failed to search users", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        utils.ListItemToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) ItemAttrSearch(ctx context.Context, req *pb.SearchReq) (*pb.PaginatedItemRes, error) {
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
		zap.L().Debug("failed to search users", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        utils.ListItemToProto(res.Data),
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
		zap.L().Debug("failed to search users", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedItemRes{
		Data:        utils.ListItemToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
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

	item := utils.ItemFromProto(req)
	if err := validation.ItemValidation(item); err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, err.Error())
	}

	res, err := h.ctrl.CreateItem(ctx, item)
	if err != nil {
		c = codes.Internal
		zap.L().Debug("failed to search users", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return utils.ItemToProto(res), nil
}

func (h *Handler) GetItem(ctx context.Context, req *pb.UuidMsg) (*pb.ItemMsg, error) {

}

func (h *Handler) UpdateItem(ctx context.Context, req *pb.ItemWithUid) (*pb.ItemMsg, error) {

}

func (h *Handler) DeleteItem(ctx context.Context, req *pb.UuidMsg) (*pb.ItemMsg, error) {

}

func (h *Handler) ListRelatedItems(ctx context.Context, req *pb.UuidMsg) (*pb.RelatedItemsList, error) {

}

func (h *Handler) ListCategoryItems(ctx context.Context, req *pb.ListCategoryItemsReq) (*pb.PaginatedItemRes, error) {

}

func (h *Handler) HitItems(ctx context.Context, req *pb.ListReq) (*pb.PaginatedItemRes, error) {

}

func (h *Handler) RecItems(ctx context.Context, req *pb.ListReq) (*pb.PaginatedItemRes, error) {

}
