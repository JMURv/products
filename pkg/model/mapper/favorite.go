package mapper

import (
	pb "github.com/JMURv/par-pro/products/api/pb"
	md "github.com/JMURv/par-pro/products/pkg/model"
)

func ListFavoriteToProto(u []*md.Favorite) []*pb.FavoriteMsg {
	res := make([]*pb.FavoriteMsg, len(u))
	for i, v := range u {
		res[i] = FavoriteToProto(v)
	}
	return res
}

func FavoriteToProto(req *md.Favorite) *pb.FavoriteMsg {
	return &pb.FavoriteMsg{
		Id:     req.ID,
		UserId: req.UserID.String(),
		ItemId: req.ItemID.String(),
		Item:   ItemToProto(&req.Item),
	}
}
