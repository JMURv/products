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

func (h *Handler) ListOrders(ctx context.Context, req *pb.ListReq) (*pb.PaginatedOrderRes, error) {
	s, c := time.Now(), codes.OK
	const op = "order.ListOrders.handler"
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

	// TODO: Create filters
	filters := make(map[string]any)
	res, err := h.ctrl.ListOrders(ctx, int(page), int(size), filters, "")
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedOrderRes{
		Data:        mapper.ListOrdersToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) ListUserOrders(ctx context.Context, req *pb.ListReq) (*pb.PaginatedOrderRes, error) {
	s, c := time.Now(), codes.OK
	const op = "order.ListUserOrders.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	uidStr, ok := ctx.Value("uid").(string)
	if !ok {
		c = codes.Unauthenticated
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrUnauthenticated.Error())
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrParseUUID.Error())
	}

	page, size := req.Page, req.Size
	if page == 0 || size == 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ListUserOrders(ctx, uid, int(page), int(size))
	if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.PaginatedOrderRes{
		Data:        mapper.ListOrdersToProto(res.Data),
		Count:       res.Count,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		HasNextPage: res.HasNextPage,
	}, nil
}

func (h *Handler) GetOrder(ctx context.Context, req *pb.Uint64Msg) (*pb.OrderMsg, error) {
	s, c := time.Now(), codes.OK
	const op = "order.GetOrder.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	orderID := req.Value
	if orderID == 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.GetOrder(ctx, orderID)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return mapper.OrderToProto(res), nil
}

func (h *Handler) CreateOrder(ctx context.Context, req *pb.OrderMsg) (*pb.Uint64Msg, error) {
	var err error
	s, c := time.Now(), codes.OK
	const op = "order.CreateOrder.handler"
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

	obj := mapper.OrderFromProto(req)
	if err = validation.Order(obj); err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, err.Error())
	}

	uid := uuid.Nil
	if uidStr, ok := ctx.Value("uid").(string); !ok {
		genPass := uuid.NewString()
		userID, err := h.sso.CreateUser(ctx, obj.FIO, obj.Email, genPass)
		if err != nil {
			zap.L().Debug("Error create user", zap.Error(err), zap.String("op", op))
			c = codes.InvalidArgument
			return nil, nil
		}

		uid, err = uuid.Parse(userID)
		if err != nil {
			zap.L().Debug("Error parse user UUID", zap.Error(err), zap.String("op", op))
			c = codes.Internal
			return nil, nil
		}
	} else {
		uid, err = uuid.Parse(uidStr)
		if err != nil {
			c = codes.InvalidArgument
			zap.L().Debug("failed to decode request", zap.String("op", op))
			return nil, status.Errorf(c, ctrl.ErrParseUUID.Error())
		}
	}

	res, err := h.ctrl.CreateOrder(ctx, uid, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.Uint64Msg{Value: res}, nil
}

func (h *Handler) UpdateOrder(ctx context.Context, req *pb.OrderMsg) (*pb.Empty, error) {
	s, c := time.Now(), codes.OK
	const op = "order.UpdateOrder.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Id == 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}
	obj := mapper.OrderFromProto(req)
	if err := validation.Order(obj); err != nil {
		c = codes.InvalidArgument
		zap.L().Debug("failed to validate request", zap.String("op", op), zap.Error(err))
		return nil, status.Errorf(c, err.Error())
	}

	err := h.ctrl.UpdateOrder(ctx, req.Id, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.Empty{}, nil
}

func (h *Handler) CancelOrder(ctx context.Context, req *pb.Uint64Msg) (*pb.Empty, error) {
	s, c := time.Now(), codes.OK
	const op = "order.CancelOrder.handler"
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Value == 0 {
		c = codes.InvalidArgument
		zap.L().Debug("failed to decode request", zap.String("op", op))
		return nil, status.Errorf(c, ctrl.ErrDecodeRequest.Error())
	}

	err := h.ctrl.CancelOrder(ctx, req.Value)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		c = codes.Internal
		return nil, status.Errorf(c, ctrl.ErrInternalError.Error())
	}

	return &pb.Empty{}, nil
}
