package etc

import (
	"context"
	"github.com/JMURv/par-pro/products/pkg/model/etc"
)

type bannerName string

const (
	Category bannerName = "category"
	Promo    bannerName = "promo"
)

func (s bannerName) String() string {
	return string(s)
}

type EtcCtrl interface {
	CreateBanner(ctx context.Context, name, pk string, req *etc.Banner) error
	UpdateBanner(ctx context.Context, name, pk string, req *etc.Banner) error
	DeleteBanner(ctx context.Context, name, pk string) error
}
