package models

import "time"

type Page struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Href  string `json:"href"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
