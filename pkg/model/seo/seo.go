package seo

import "time"

type SEO struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`

	OGTitle       string `json:"OGTitle"`
	OGDescription string `json:"OGDescription"`
	OGImage       string `json:"OGImage"`

	OBJName string `json:"obj_name"`
	OBJPK   string `json:"obj_pk"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
