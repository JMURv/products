package model

import (
	"github.com/lib/pq"
	"time"
)

type Category struct {
	Slug            string `json:"slug" gorm:"primaryKey;unique;not null"`
	Title           string `json:"title" gorm:"type:varchar(255);unique;not null"`
	ProductQuantity int    `json:"product_quantity"`
	Src             string `json:"src" gorm:"type:varchar(255)"`
	Alt             string `json:"alt" gorm:"type:varchar(255)"`

	ParentSlug     *string    `json:"parent_slug" gorm:"type:varchar(255)"`
	ParentCategory *Category  `json:"parent_category" gorm:"foreignKey:ParentSlug;constraint:OnDelete:SET NULL"`
	Children       []Category `json:"children" gorm:"foreignKey:ParentSlug;constraint:OnDelete:SET NULL"`

	Items   []*Item  `json:"items" gorm:"many2many:item_categories;constraint:OnDelete:SET NULL"`
	Filters []Filter `json:"filters" gorm:"foreignKey:CategorySlug;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Filter struct {
	ID     uint64         `json:"id" gorm:"primaryKey"`
	Name   string         `json:"name" gorm:"type:varchar(255)"`
	Values pq.StringArray `json:"values" gorm:"type:varchar(255)[]"`

	FilterType string   `json:"filter_type" gorm:"type:varchar(50)"` // "equality", "range"
	MinValue   *float64 `json:"min_value" gorm:""`                   // For range filters, specify the minimum value (e.g., min price)
	MaxValue   *float64 `json:"max_value" gorm:""`                   // For range filters, specify the maximum value (e.g., max price)

	CategorySlug string `json:"category_slug" gorm:"type:varchar(255)"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
