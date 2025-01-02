package model

import (
	"github.com/google/uuid"
	"time"
)

const OrderStatusPending = "pending"
const OrderStatusConfirmed = "confirmed"
const OrderStatusDelivered = "delivered"
const OrderStatusCancelled = "cancelled"

type DeliveryType string

const (
	DeliveryTypePickup   DeliveryType = "самовывоз"
	DeliveryTypeDelivery DeliveryType = "доставка"
)

type PaymentMethod string

const (
	PaymentMethodCard  PaymentMethod = "наличными"
	PaymentMethodCash  PaymentMethod = "банковской картой"
	PaymentMethodOther PaymentMethod = "расчетный счет компании"
)

type PaginatedOrderData struct {
	Data        []*Order `json:"data"`
	Count       int64    `json:"count"`
	TotalPages  int      `json:"total_pages"`
	CurrentPage int      `json:"current_page"`
	HasNextPage bool     `json:"has_next_page"`
}

type Order struct {
	ID          uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Status      string  `json:"status" gorm:"not null"`
	TotalAmount float64 `json:"total_amount" gorm:"not null"`

	FIO           string `json:"fio" gorm:"type:varchar(255);not null"`
	Tel           string `json:"tel" gorm:"type:varchar(20);not null"`
	Email         string `json:"email" gorm:"type:varchar(255);not null"`
	Address       string `json:"address" gorm:"type:varchar(255);not null"`
	Delivery      string `json:"delivery" gorm:"type:varchar(255);not null"`
	PaymentMethod string `json:"payment_method" gorm:"type:varchar(255);not null"`

	UserID     uuid.UUID   `json:"user_id" gorm:"not null"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Quantity int    `json:"quantity" gorm:"not null"`

	OrderID uint64    `json:"order_id"`
	ItemID  uuid.UUID `json:"item_id"`
	Item    Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
