package seo

import (
	"context"
	"github.com/JMURv/par-pro/products/pkg/model/seo"
)

type seoName string

const (
	Category seoName = "category"
	Item     seoName = "item"
	Promo    seoName = "promo"
)

func (s seoName) String() string {
	return string(s)
}

type SEOCtrl interface {
	CreateSEO(ctx context.Context, name, pk string, seo *seo.SEO) error
	UpdateSEO(ctx context.Context, name, pk string, seo *seo.SEO) error
	DeleteSEO(ctx context.Context, name, pk string) error
}
