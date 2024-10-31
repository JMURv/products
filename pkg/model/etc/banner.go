package etc

import "time"

type Banner struct {
	ID      uint64 `json:"id"`
	OBJName string `json:"obj_name"`
	OBJPK   string `json:"obj_pk"`

	Slides    []BannerSlide `json:"slides"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type BannerSlide struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Src         string `json:"src"`
	Alt         string `json:"alt"`
	ButtonText  string `json:"button_text"`
	ButtonHref  string `json:"button_href"`

	BannerID uint64 `json:"banner_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
