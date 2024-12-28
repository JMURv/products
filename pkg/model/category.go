package model

import (
	"time"
)

type PaginatedCategoryData struct {
	Data        []*Category `json:"data"`
	Count       int64       `json:"count"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	HasNextPage bool        `json:"has_next_page"`
}

type Category struct {
	ID              uint64 `json:"id"`
	Slug            string `json:"slug"`
	Title           string `json:"title"`
	ProductQuantity int64  `json:"product_quantity"`
	Src             string `json:"src"`
	Alt             string `json:"alt"`

	ParentSlug     string     `json:"parent_slug"`
	ParentCategory *Category  `json:"parent_category"`
	Children       []Category `json:"children"`

	Filters []Filter `json:"filters"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedFilterData struct {
	Data        []*Filter `json:"data"`
	Count       int64     `json:"count"`
	TotalPages  int       `json:"total_pages"`
	CurrentPage int       `json:"current_page"`
	HasNextPage bool      `json:"has_next_page"`
}

type Filter struct {
	ID     uint64   `json:"id"`
	Name   string   `json:"name"`
	Values []string `json:"values"`

	FilterType string  `json:"filter_type"` // "equality", "range"
	MinValue   float64 `json:"min_value"`   // For range filters, specify the minimum value (e.g., min price)
	MaxValue   float64 `json:"max_value"`   // For range filters, specify the maximum value (e.g., max price)

	CategorySlug string `json:"category_slug"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
