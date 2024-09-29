package etc

import "time"

type Banner struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	OBJName string `json:"obj_name"`
	OBJPK   string `json:"obj_pk"`

	Slides    []BannerSlide `json:"slides" gorm:"foreignKey:BannerID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type BannerSlide struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null;type:varchar(255)"`
	Description string `json:"description"`
	Src         string `json:"src" gorm:"not null;type:varchar(255)"`
	Alt         string `json:"alt" gorm:"type:varchar(255)"`
	ButtonText  string `json:"button_text" gorm:"type:varchar(255)"`
	ButtonHref  string `json:"button_href" gorm:"type:varchar(255)"`

	BannerID uint64 `json:"banner_id"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
