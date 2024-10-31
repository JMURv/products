package model

import (
	"github.com/google/uuid"
	"time"
)

type Cart struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	UserID uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;unique"`
	Items  []CartItem `json:"items" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CartItem struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Quantity int       `json:"quantity" gorm:"type:integer;not null"`

	CartID uuid.UUID `json:"cart_id" gorm:"type:uuid;not null"`
	ItemID uuid.UUID `json:"item_id" gorm:"type:uuid;not null"`
	Item   Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
