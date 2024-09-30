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
	Slug           string           `json:"slug" gorm:"primaryKey;unique;not null;type:varchar(255)"`
	Title          string           `json:"title" gorm:"type:varchar(255);not null"`
	Description    string           `json:"description" gorm:"type:text;not null"`
	Src            string           `json:"src" gorm:"type:varchar(255);not null"`
	Alt            string           `json:"alt" gorm:"type:varchar(255)"`
	LastsTo        time.Time        `json:"lasts_to" gorm:"not null"`
	PromotionItems []*PromotionItem `json:"promotion_items" gorm:"foreignKey:PromotionSlug;constraint:OnDelete:CASCADE"`

	Banner    *etc.Banner `json:"banner" gorm:"-"`
	SEO       *seo.SEO    `json:"seo" gorm:"-"`
	CreatedAt time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type PaginatedPromoItemsData struct {
	Data        []*PromotionItem `json:"data"`
	Count       int64            `json:"count"`
	TotalPages  int              `json:"total_pages"`
	CurrentPage int              `json:"current_page"`
	HasNextPage bool             `json:"has_next_page"`
}

type PromotionItem struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	Discount      int    `json:"discount"`
	PromotionSlug string `json:"promotion_slug" gorm:"type:varchar(255)"`

	ItemID uuid.UUID `json:"item_id" gorm:"type:uuid"`
	Item   Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
