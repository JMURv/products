package model

import (
	"github.com/JMURv/par-pro/products/pkg/utils/slugify"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
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

func MustPrecreatePromotion(conn *gorm.DB) {
	var count int64
	if err := conn.Model(&Promotion{}).Count(&count).Error; err != nil {
		panic(err)
	}

	if count == 0 {
		bytes, err := os.ReadFile("db/precreate/promotions.json")
		if err != nil {
			panic(err)
		}

		promos := make([]*Promotion, 0, 4)
		if err = json.Unmarshal(bytes, &promos); err != nil {
			panic(err)
		}

		for _, v := range promos {
			slug := slugify.Slugify(v.Title)
			promo := &Promotion{
				Slug:           slug,
				Title:          v.Title,
				Description:    v.Description,
				Src:            v.Src,
				Alt:            v.Alt,
				LastsTo:        time.Now().AddDate(0, 0, 7),
				PromotionItems: v.PromotionItems,
			}
			if err = conn.Create(promo).Error; err != nil {
				panic(err)
			}
		}

		zap.L().Debug("Promotions have been created")
	} else {
		zap.L().Debug("Promotions already exist")
	}
}
