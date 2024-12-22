package model

import (
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
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	QuantityInStock int       `json:"quantity_in_stock"`
	InStock         bool      `json:"in_stock"`
	Src             string    `json:"src"`
	Alt             string    `json:"alt"`
	IsHit           bool      `json:"is_hit"`
	IsRec           bool      `json:"is_rec"`
	Article         string    `json:"article"`

	Categories   []Category `json:"categories"`
	ParentItemID uuid.UUID  `json:"parent_item_id"`

	Media           []ItemMedia      `json:"media"`
	Attributes      []ItemAttribute  `json:"attributes"`
	Variants        []Item           `json:"variants"`
	RelatedProducts []RelatedProduct `json:"related_products"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedItemAttrData struct {
	Data        []*ItemAttribute `json:"data"`
	Count       int64            `json:"count"`
	TotalPages  int              `json:"total_pages"`
	CurrentPage int              `json:"current_page"`
	HasNextPage bool             `json:"has_next_page"`
}

type ItemAttribute struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`

	ItemID uuid.UUID `json:"item_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemMedia struct {
	ID  uint64 `json:"id"`
	Src string `json:"src"`
	Alt string `json:"alt"`

	ItemID uuid.UUID `json:"item_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemCategory struct {
	ItemID       uuid.UUID `json:"item_id"`
	CategorySlug string    `json:"category_slug"`
}

type RelatedProduct struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	ItemID        uuid.UUID `json:"item_id"`
	RelatedItemID uuid.UUID `json:"related_item_id"`
	RelatedItem   Item      `json:"related_item"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
