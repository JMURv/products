package model

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"time"
)

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

	Media           []ItemMedia      `json:"media" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Attributes      []ItemAttribute  `json:"attributes" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Variants        []Item           `json:"variants" gorm:"foreignKey:ParentItemID;constraint:OnDelete:CASCADE"`
	RelatedProducts []RelatedProduct `json:"related_products" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
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

func MustPrecreateItems(conn *gorm.DB) {
	var count int64
	if err := conn.Model(&Item{}).Count(&count).Error; err != nil {
		panic(err)
	}

	if count == 0 {
		bytes, err := os.ReadFile("db/precreate/items.json")
		if err != nil {
			panic(err)
		}

		items := make([]*Item, 0, 2)
		if err = json.Unmarshal(bytes, &items); err != nil {
			panic(err)
		}

		for _, item := range items {
			i := &Item{
				ID:              item.ID,
				Title:           item.Title,
				Description:     item.Description,
				Price:           item.Price,
				InStock:         item.InStock,
				Src:             item.Src,
				Alt:             item.Alt,
				IsHit:           item.IsHit,
				IsRec:           item.IsRec,
				Media:           item.Media,
				Attributes:      item.Attributes,
				Variants:        item.Variants,
				RelatedProducts: item.RelatedProducts,
			}

			for _, category := range item.Categories {
				var c Category
				if err := conn.Where("slug = ?", category.Slug).First(&c).Error; err == nil {
					i.Categories = append(i.Categories, &c)
				}
			}

			if err = conn.Create(i).Error; err != nil {
				panic(err)
			}
		}

		zap.L().Debug("Items have been created")
	} else {
		zap.L().Debug("Items already exist")
	}
}
