package mapper

import (
	pb "github.com/JMURv/par-pro/products/api/pb"
	md "github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ListOrdersToProto(req []*md.Order) []*pb.OrderMsg {
	orders := make([]*pb.OrderMsg, len(req))
	for i := 0; i < len(req); i++ {
		orders[i] = OrderToProto(req[i])
	}

	return orders
}

func ListOrderItemsToProto(req []*md.OrderItem) []*pb.OrderItem {
	items := make([]*pb.OrderItem, len(req))
	for i := 0; i < len(req); i++ {
		items[i] = OrderItemToProto(req[i])
	}
	return items
}

func OrderItemToProto(req *md.OrderItem) *pb.OrderItem {
	return &pb.OrderItem{
		Id:        req.ID,
		Quantity:  uint32(req.Quantity),
		OrderId:   req.OrderID,
		Item:      ItemToProto(&req.Item),
		CreatedAt: timestamppb.New(req.CreatedAt),
		UpdatedAt: timestamppb.New(req.UpdatedAt),
	}

}

func OrderToProto(req *md.Order) *pb.OrderMsg {
	order := &pb.OrderMsg{
		Id:            req.ID,
		Status:        req.Status,
		Total:         float32(req.TotalAmount),
		Fio:           req.FIO,
		Tel:           req.Tel,
		Email:         req.Email,
		Address:       req.Address,
		Delivery:      req.Delivery,
		PaymentMethod: req.PaymentMethod,
		UserId:        req.UserID.String(),
		CreatedAt:     timestamppb.New(req.CreatedAt),
		UpdatedAt:     timestamppb.New(req.UpdatedAt),
	}
	if len(req.OrderItems) > 0 {
		order.Items = ListOrderItemsToProto(req.OrderItems)
	}

	return order
}

func ListOrderItemsFromProto(req []*pb.OrderItem) []*md.OrderItem {
	items := make([]*md.OrderItem, len(req))
	for i := 0; i < len(req); i++ {
		items[i] = OrderItemFromProto(req[i])
	}
	return items
}

func OrderItemFromProto(req *pb.OrderItem) *md.OrderItem {
	return &md.OrderItem{
		ID:        req.Id,
		Quantity:  int(req.Quantity),
		OrderID:   req.OrderId,
		Item:      *ItemFromProto(req.Item),
		CreatedAt: req.CreatedAt.AsTime(),
		UpdatedAt: req.UpdatedAt.AsTime(),
	}
}

func OrderFromProto(req *pb.OrderMsg) *md.Order {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		zap.L().Debug("failed to parse user id")
		return nil
	}

	order := &md.Order{
		ID:            req.Id,
		Status:        req.Status,
		TotalAmount:   float64(req.Total),
		FIO:           req.Fio,
		Tel:           req.Tel,
		Email:         req.Email,
		Address:       req.Address,
		Delivery:      req.Delivery,
		PaymentMethod: req.PaymentMethod,
		UserID:        uid,
		CreatedAt:     req.CreatedAt.AsTime(),
		UpdatedAt:     req.UpdatedAt.AsTime(),
	}
	if len(req.Items) > 0 {
		order.OrderItems = ListOrderItemsFromProto(req.Items)
	}

	return order
}
