package models

type Contact struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	ImageURL string `json:"image_url"`
}
