package model

import (
	"github.com/google/uuid"
)

type Favorite struct {
	ID uint64 `json:"id"`

	UserID uuid.UUID `json:"user_id"`
	ItemID uuid.UUID `json:"item_id"`
	Item   Item      `json:"item" gorm:"constraint:OnDelete:CASCADE"`
}
