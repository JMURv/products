package seo

import "time"

type SEO struct {
	ID          uint64 `json:"id" gorm:"primaryKey;index"`
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"not null"`
	Keywords    string `json:"keywords" gorm:"not null"`

	OGTitle       string `json:"OGTitle" gorm:"type:varchar(255)"`
	OGDescription string `json:"OGDescription"`
	OGImage       string `json:"OGImage"`

	OBJName string `json:"obj_name"`
	OBJPK   string `json:"obj_pk"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
