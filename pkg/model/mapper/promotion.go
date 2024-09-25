package mapper

import (
	pb "github.com/JMURv/par-pro/products/api/pb"
	md "github.com/JMURv/par-pro/products/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ListPromosToProto(req []*md.Promotion) []*pb.PromoMsg {
	res := make([]*pb.PromoMsg, len(req))
	for i, v := range req {
		res[i] = PromoToProto(v)
	}
	return res
}

func PromoToProto(req *md.Promotion) *pb.PromoMsg {
	res := &pb.PromoMsg{
		Slug:        req.Slug,
		Title:       req.Title,
		Description: req.Description,
		Src:         req.Src,
		Alt:         req.Alt,
		LastsTo: &timestamppb.Timestamp{
			Seconds: req.LastsTo.Unix(),
			Nanos:   int32(req.LastsTo.Nanosecond()),
		},
		CreatedAt: &timestamppb.Timestamp{
			Seconds: req.CreatedAt.Unix(),
			Nanos:   int32(req.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: req.UpdatedAt.Unix(),
			Nanos:   int32(req.UpdatedAt.Nanosecond()),
		},
	}

	return res
}

func PromoItemToProto(req *md.PromotionItem) *pb.PromoItem {
	return &pb.PromoItem{
		Id:            req.ID,
		Discount:      uint32(req.Discount),
		PromotionSlug: req.PromotionSlug,
		ItemId:        req.ItemID.String(),
		Item:          ItemToProto(&req.Item),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: req.CreatedAt.Unix(),
			Nanos:   int32(req.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: req.UpdatedAt.Unix(),
			Nanos:   int32(req.UpdatedAt.Nanosecond()),
		},
	}
}

func ListPromoItemsToProto(req []*md.PromotionItem) []*pb.PromoItem {
	res := make([]*pb.PromoItem, len(req))
	for i, v := range req {
		res[i] = PromoItemToProto(v)
	}
	return res
}

func PromoFromProto(req *pb.PromoMsg) *md.Promotion {
	return &md.Promotion{
		Slug:        req.Slug,
		Title:       req.Title,
		Description: req.Description,
		Src:         req.Src,
		Alt:         req.Alt,
		LastsTo:     req.LastsTo.AsTime(),
		CreatedAt:   req.CreatedAt.AsTime(),
		UpdatedAt:   req.UpdatedAt.AsTime(),
	}
}
