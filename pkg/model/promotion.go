package model

import (
	"github.com/google/uuid"
	"time"
)

type Promotion struct {
	Slug           string           `json:"slug" gorm:"primaryKey;unique;not null;type:varchar(255)"`
	Title          string           `json:"title" gorm:"type:varchar(255);not null"`
	Description    string           `json:"description" gorm:"type:text;not null"`
	Src            string           `json:"src" gorm:"type:varchar(255);not null"`
	Alt            string           `json:"alt" gorm:"type:varchar(255)"`
	LastsTo        time.Time        `json:"lasts_to" gorm:"not null"`
	PromotionItems []*PromotionItem `json:"promotion_items" gorm:"foreignKey:PromotionSlug;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
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
