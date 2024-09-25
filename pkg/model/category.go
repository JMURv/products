package model

import (
	"github.com/JMURv/par-pro/pkg/consts"
	"github.com/JMURv/par-pro/pkg/utils/slugify"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
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

	SEO     SEO      `json:"seo" gorm:"foreignKey:CategorySlug;constraint:OnDelete:CASCADE"`
	Items   []*Item  `json:"items" gorm:"many2many:item_categories;constraint:OnDelete:SET NULL"`
	Filters []Filter `json:"filters" gorm:"foreignKey:CategorySlug;constraint:OnDelete:CASCADE"`
	Banner  Banner   `json:"banner" gorm:"foreignKey:CategorySlug;constraint:OnDelete:SET NULL"`

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

func MustPrecreateCategory(conn *gorm.DB) {
	var count int64
	if err := conn.Model(&Category{}).Count(&count).Error; err != nil {
		panic(err)
	}

	if count == 0 {
		bytes, err := os.ReadFile("db/precreate/category.json")
		if err != nil {
			panic(err)
		}

		p := make([]*Category, 0, 5)
		if err = json.Unmarshal(bytes, &p); err != nil {
			panic(err)
		}

		for _, category := range p {
			slug := slugify.Slugify(category.Title)
			c := &Category{
				Slug:            slug,
				Title:           category.Title,
				ProductQuantity: category.ProductQuantity,
				Src:             category.Src,
				Alt:             category.Alt,
				ParentSlug:      category.ParentSlug,
				Filters:         category.Filters,
				Banner: Banner{
					CategorySlug: &slug,
				},
				SEO: SEO{
					Title:         category.SEO.Title,
					Description:   category.SEO.Description,
					Keywords:      category.SEO.Keywords,
					OGTitle:       category.SEO.OGTitle,
					OGDescription: category.SEO.OGDescription,
					OGImage:       consts.DefaultImagePath,
					CategorySlug:  &slug,
				},
			}
			if err = conn.Create(c).Error; err != nil {
				panic(err)
			}
		}

		zap.L().Debug("Categories have been created")
	} else {
		zap.L().Debug("Categories already exist")
	}
}
