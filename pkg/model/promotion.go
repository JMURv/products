package model

import (
	"github.com/JMURv/par-pro/products/pkg/model/etc"
	"github.com/JMURv/par-pro/products/pkg/model/seo"
	"github.com/google/uuid"
	"time"
)

type PaginatedPromosData struct {
	Data        []*Promotion `json:"data"`
	Count       int64        `json:"count"`
	TotalPages  int          `json:"total_pages"`
	CurrentPage int          `json:"current_page"`
	HasNextPage bool         `json:"has_next_page"`
}

type Promotion struct {
	Slug           string           `json:"slug"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Src            string           `json:"src"`
	Alt            string           `json:"alt"`
	LastsTo        time.Time        `json:"lasts_to"`
	PromotionItems []*PromotionItem `json:"promotion_items"`

	Banner    *etc.Banner `json:"banner"`
	SEO       *seo.SEO    `json:"seo"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type PaginatedPromoItemsData struct {
	Data        []*PromotionItem `json:"data"`
	Count       int64            `json:"count"`
	TotalPages  int              `json:"total_pages"`
	CurrentPage int              `json:"current_page"`
	HasNextPage bool             `json:"has_next_page"`
}

type PromotionItem struct {
	ID            uint64    `json:"id"`
	Discount      int       `json:"discount"`
	PromotionSlug string    `json:"promotion_slug"`
	ItemID        uuid.UUID `json:"item_id"`
	Item          Item      `json:"item"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
