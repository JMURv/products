package mapper

import (
	pb "github.com/JMURv/par-pro/products/api/pb"
	md "github.com/JMURv/par-pro/products/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ListCategoryToProto(req []*md.Category) []*pb.CategoryMsg {
	res := make([]*pb.CategoryMsg, len(req))
	for i, v := range req {
		res[i] = CategoryToProto(v)
	}
	return res
}

func CategoryToProto(req *md.Category) *pb.CategoryMsg {
	res := &pb.CategoryMsg{
		Slug:               req.Slug,
		Title:              req.Title,
		ProductQuantity:    uint64(req.ProductQuantity),
		Src:                req.Src,
		Alt:                req.Alt,
		ParentSlug:         req.ParentSlug,
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
			res.Items = append(res.Items, ItemToProto(&v))
		}
	}

	if len(req.Filters) > 0 {
		for _, v := range req.Filters {
			res.Filters = append(
				res.Filters, &pb.Filter{
					Id:           v.ID,
					Name:         v.Name,
					Values:       v.Values,
					CategorySlug: v.CategorySlug,
					FilterType:   v.FilterType,
					MinValue:     float32(v.MinValue),
					MaxValue:     float32(v.MaxValue),
					CreatedAt: &timestamppb.Timestamp{
						Seconds: v.CreatedAt.Unix(),
						Nanos:   int32(v.CreatedAt.Nanosecond()),
					},
					UpdatedAt: &timestamppb.Timestamp{
						Seconds: v.UpdatedAt.Unix(),
						Nanos:   int32(v.UpdatedAt.Nanosecond()),
					},
				},
			)
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
		ParentSlug:      req.ParentSlug,
		CreatedAt:       req.CreatedAt.AsTime(),
		UpdatedAt:       req.UpdatedAt.AsTime(),
	}

	if req.Parent_CategoryMsg != nil {
		res.ParentCategory = CategoryFromProto(req.Parent_CategoryMsg)
	}

	if len(req.Children) > 0 {
		for _, v := range req.Children {
			res.Children = append(res.Children, *CategoryFromProto(v))
		}
	}

	if len(req.Items) > 0 {
		for _, v := range req.Items {
			itm := ItemFromProto(v)
			res.Items = append(res.Items, *itm)
		}
	}

	if len(req.Filters) > 0 {
		for _, v := range req.Filters {
			minVal := float64(v.MinValue)
			maxVal := float64(v.MaxValue)
			res.Filters = append(
				res.Filters, md.Filter{
					ID:           v.Id,
					Name:         v.Name,
					Values:       v.Values,
					CategorySlug: v.CategorySlug,
					FilterType:   v.FilterType,
					MinValue:     minVal,
					MaxValue:     maxVal,
					CreatedAt:    v.CreatedAt.AsTime(),
					UpdatedAt:    v.UpdatedAt.AsTime(),
				},
			)
		}
	}

	return res
}

func ListFiltersToProto(req []*md.Filter) []*pb.Filter {
	res := make([]*pb.Filter, len(req))
	for i, v := range req {
		minVal := float32(v.MinValue)
		maxVal := float32(v.MaxValue)
		res[i] = &pb.Filter{
			Id:           v.ID,
			Name:         v.Name,
			Values:       v.Values,
			CategorySlug: v.CategorySlug,
			FilterType:   v.FilterType,
			MinValue:     minVal,
			MaxValue:     maxVal,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: v.CreatedAt.Unix(),
				Nanos:   int32(v.CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: v.UpdatedAt.Unix(),
				Nanos:   int32(v.UpdatedAt.Nanosecond()),
			},
		}
	}
	return res
}
