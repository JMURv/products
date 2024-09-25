package utils

import (
	pb "github.com/JMURv/par-pro/products/api/pb"
	md "github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ListItemToProto(u []*md.Item) []*pb.ItemMsg {
	res := make([]*pb.ItemMsg, len(u))
	for i, v := range u {
		res[i] = ItemToProto(v)
	}
	return res
}

func ItemToProto(req *md.Item) *pb.ItemMsg {
	item := &pb.ItemMsg{
		Id:              req.ID.String(),
		Article:         req.Article,
		Title:           req.Title,
		Description:     req.Description,
		Price:           float32(req.Price),
		QuantityInStock: uint32(req.QuantityInStock),
		Src:             req.Src,
		Alt:             req.Alt,
		InStock:         req.InStock,
		IsHit:           req.IsHit,
		IsRec:           req.IsRec,
		ParentItemId:    req.ParentItemID.String(),
		Media:           ListItemMediaToProto(req.Media),
		Attributes:      ListItemAttributesToProto(req.Attributes),
		RelatedProducts: ListRelatedProductsToProto(req.RelatedProducts),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: req.CreatedAt.Unix(),
			Nanos:   int32(req.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: req.UpdatedAt.Unix(),
			Nanos:   int32(req.UpdatedAt.Nanosecond()),
		},
	}

	if len(req.Categories) > 0 {
		res := make([]*pb.CategoryMsg, len(req.Categories))
		for i, v := range req.Categories {
			res[i] = CategoryToProto(v)
		}
		item.Categories = res
	}

	if len(req.Variants) > 0 {
		res := make([]*pb.ItemMsg, len(req.Variants))
		for i, v := range req.Variants {
			res[i] = ItemToProto(&v)
		}
		item.Variants = res
	}

	return item
}

func ListRelatedProductsToProto(u []md.RelatedProduct) []*pb.RelatedProduct {
	res := make([]*pb.RelatedProduct, len(u))
	for i, v := range u {
		res[i] = RelatedProductsToProto(&v)
	}
	return res
}

func RelatedProductsToProto(req *md.RelatedProduct) *pb.RelatedProduct {
	return &pb.RelatedProduct{
		Id:            req.ID,
		ItemId:        req.ItemID.String(),
		RelatedItemId: req.RelatedItemID.String(),
		RelatedItem:   ItemToProto(&req.RelatedItem),
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

func ListItemAttributesToProto(u []md.ItemAttribute) []*pb.ItemAttribute {
	res := make([]*pb.ItemAttribute, len(u))
	for i, v := range u {
		res[i] = ItemAttributesToProto(&v)
	}
	return res
}

func ItemAttributesToProto(req *md.ItemAttribute) *pb.ItemAttribute {
	return &pb.ItemAttribute{
		Id:     req.ID,
		Name:   req.Name,
		Value:  req.Value,
		ItemId: req.ItemID.String(),
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

func ListItemMediaToProto(u []md.ItemMedia) []*pb.ItemMedia {
	res := make([]*pb.ItemMedia, len(u))
	for i, v := range u {
		res[i] = ItemMediaToProto(&v)
	}
	return res
}

func ItemMediaToProto(req *md.ItemMedia) *pb.ItemMedia {
	return &pb.ItemMedia{
		Id:     req.ID,
		Src:    req.Src,
		Alt:    req.Alt,
		ItemId: req.ItemID.String(),
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

func ItemFromProto(req *pb.ItemMsg) *md.Item {
	modelItem := &md.Item{
		Article:         req.Article,
		Title:           req.Title,
		Description:     req.Description,
		Price:           float64(req.Price),
		QuantityInStock: int(req.QuantityInStock),
		Src:             req.Src,
		Alt:             req.Alt,
		InStock:         req.InStock,
		IsHit:           req.IsHit,
		IsRec:           req.IsRec,
		CreatedAt:       req.CreatedAt.AsTime(),
		UpdatedAt:       req.UpdatedAt.AsTime(),
	}

	uid, err := uuid.Parse(req.Id)
	if err != nil {
		zap.L().Debug("failed to parse user id")
	} else {
		modelItem.ID = uid
	}

	parentItemID, err := uuid.Parse(req.ParentItemId)
	if err != nil {
		zap.L().Debug("failed to parse parent item id")
	} else {
		modelItem.ParentItemID = &parentItemID
	}

	if len(req.Categories) > 0 {
		res := make([]*md.Category, len(req.Categories))
		for i, v := range req.Categories {
			res[i] = CategoryFromProto(v)
		}
		modelItem.Categories = res
	}

	if len(req.Media) > 0 {
		res := make([]md.ItemMedia, len(req.Media))
		for i, v := range req.Media {
			iMedia := md.ItemMedia{
				ID:        v.Id,
				Src:       v.Src,
				Alt:       v.Alt,
				CreatedAt: v.CreatedAt.AsTime(),
				UpdatedAt: v.UpdatedAt.AsTime(),
			}

			uid, err := uuid.Parse(v.ItemId)
			if err != nil {
				zap.L().Debug("failed to parse user id")
			} else {
				iMedia.ItemID = uid
			}

			res[i] = iMedia
		}
		modelItem.Media = res
	}

	if len(req.Attributes) > 0 {
		res := make([]md.ItemAttribute, len(req.Attributes))
		for i, v := range req.Attributes {
			res[i] = md.ItemAttribute{
				ID:        v.Id,
				Name:      v.Name,
				Value:     v.Value,
				CreatedAt: v.CreatedAt.AsTime(),
				UpdatedAt: v.UpdatedAt.AsTime(),
			}
		}
		modelItem.Attributes = res
	}

	if len(req.RelatedProducts) > 0 {
		res := make([]md.RelatedProduct, len(req.RelatedProducts))
		for i, v := range req.RelatedProducts {
			pr := md.RelatedProduct{
				ID:          v.Id,
				RelatedItem: *ItemFromProto(v.RelatedItem),
				CreatedAt:   v.CreatedAt.AsTime(),
				UpdatedAt:   v.UpdatedAt.AsTime(),
			}

			itemUid, err := uuid.Parse(v.ItemId)
			if err != nil {
				zap.L().Debug("failed to parse user id")
			} else {
				pr.ItemID = itemUid
			}

			relatedItemUid, err := uuid.Parse(v.RelatedItemId)
			if err != nil {
				zap.L().Debug("failed to parse user id")
			} else {
				pr.RelatedItemID = relatedItemUid
			}

			res[i] = pr
		}
		modelItem.RelatedProducts = res
	}

	return modelItem
}

func CategoryToProto(req *md.Category) *pb.CategoryMsg {
	res := &pb.CategoryMsg{
		Slug:               req.Slug,
		Title:              req.Title,
		ProductQuantity:    uint64(req.ProductQuantity),
		Src:                req.Src,
		Alt:                req.Alt,
		ParentSlug:         *req.ParentSlug,
		Parent_CategoryMsg: CategoryToProto(req.ParentCategory),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: req.CreatedAt.Unix(),
			Nanos:   int32(req.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: req.UpdatedAt.Unix(),
			Nanos:   int32(req.UpdatedAt.Nanosecond()),
		},
	}

	if len(req.Children) > 0 {
		for _, v := range req.Children {
			res.Children = append(res.Children, CategoryToProto(&v))
		}
	}

	if len(req.Items) > 0 {
		for _, v := range req.Items {
			res.Items = append(res.Items, ItemToProto(v))
		}
	}

	if len(req.Filters) > 0 {
		for _, v := range req.Filters {
			res.Filters = append(res.Filters, &pb.Filter{
				Id:           v.ID,
				Name:         v.Name,
				Values:       v.Values,
				CategorySlug: v.CategorySlug,
				FilterType:   v.FilterType,
				MinValue:     float32(*v.MinValue),
				MaxValue:     float32(*v.MaxValue),
				CreatedAt: &timestamppb.Timestamp{
					Seconds: v.CreatedAt.Unix(),
					Nanos:   int32(v.CreatedAt.Nanosecond()),
				},
				UpdatedAt: &timestamppb.Timestamp{
					Seconds: v.UpdatedAt.Unix(),
					Nanos:   int32(v.UpdatedAt.Nanosecond()),
				},
			})
		}
	}

	return res
}

func CategoryFromProto(req *pb.CategoryMsg) *md.Category {
	res := &md.Category{
		Slug:            req.Slug,
		Title:           req.Title,
		ProductQuantity: int(req.ProductQuantity),
		Src:             req.Src,
		Alt:             req.Alt,
		ParentSlug:      &req.ParentSlug,
		ParentCategory:  CategoryFromProto(req.Parent_CategoryMsg),
		CreatedAt:       req.CreatedAt.AsTime(),
		UpdatedAt:       req.UpdatedAt.AsTime(),
	}

	if len(req.Children) > 0 {
		for _, v := range req.Children {
			res.Children = append(res.Children, *CategoryFromProto(v))
		}
	}

	if len(req.Items) > 0 {
		for _, v := range req.Items {
			res.Items = append(res.Items, ItemFromProto(v))
		}
	}

	if len(req.Filters) > 0 {
		for _, v := range req.Filters {
			minVal := float64(v.MinValue)
			maxVal := float64(v.MaxValue)
			res.Filters = append(res.Filters, md.Filter{
				ID:           v.Id,
				Name:         v.Name,
				Values:       v.Values,
				CategorySlug: v.CategorySlug,
				FilterType:   v.FilterType,
				MinValue:     &minVal,
				MaxValue:     &maxVal,
				CreatedAt:    v.CreatedAt.AsTime(),
				UpdatedAt:    v.UpdatedAt.AsTime(),
			})
		}
	}

	return res
}
