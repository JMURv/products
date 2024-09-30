package model

import (
	"github.com/JMURv/par-pro/products/pkg/model/seo"
	"github.com/google/uuid"
	"time"
)

type PaginatedItemsData struct {
	Data        []*Item `json:"data"`
	Count       int64   `json:"count"`
	TotalPages  int     `json:"total_pages"`
	CurrentPage int     `json:"current_page"`
	HasNextPage bool    `json:"has_next_page"`
}

type Item struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title           string    `json:"title" gorm:"type:varchar(255);not null"`
	Description     string    `json:"description" gorm:"type:text"`
	Price           float64   `json:"price"`
	QuantityInStock int       `json:"quantity_in_stock"`
	InStock         bool      `json:"in_stock" gorm:"default:true"`
	Src             string    `json:"src" gorm:"type:varchar(255)"`
	Alt             string    `json:"alt" gorm:"type:varchar(255)"`
	IsHit           bool      `json:"is_hit"`
	IsRec           bool      `json:"is_rec"`
	Article         string    `json:"article" gorm:"type:varchar(255)"`

	Categories   []*Category `json:"categories" gorm:"many2many:item_categories;constraint:OnDelete:CASCADE"`
	ParentItemID *uuid.UUID  `json:"parent_item_id" gorm:"type:uuid"`

	SEO             *seo.SEO         `json:"seo" gorm:"-"`
	Media           []ItemMedia      `json:"media" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Attributes      []ItemAttribute  `json:"attributes" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Variants        []Item           `json:"variants" gorm:"foreignKey:ParentItemID;constraint:OnDelete:CASCADE"`
	RelatedProducts []RelatedProduct `json:"related_products" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PaginatedItemAttrData struct {
	Data        []*ItemAttribute `json:"data"`
	Count       int64            `json:"count"`
	TotalPages  int              `json:"total_pages"`
	CurrentPage int              `json:"current_page"`
	HasNextPage bool             `json:"has_next_page"`
}

type ItemAttribute struct {
	ID    uint64 `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"type:varchar(255)"`
	Value string `json:"value" gorm:"type:text"`

	ItemID uuid.UUID `json:"item_id"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ItemMedia struct {
	ID  uint64 `json:"id" gorm:"primaryKey"`
	Src string `json:"src" gorm:"type:varchar(255)"`
	Alt string `json:"alt" gorm:"type:varchar(255)"`

	ItemID uuid.UUID `json:"item_id"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type RelatedProduct struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	ItemID        uuid.UUID `json:"item_id"`
	RelatedItemID uuid.UUID `json:"related_item_id"`
	RelatedItem   Item      `json:"related_item" gorm:"foreignKey:RelatedItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
