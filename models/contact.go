package models

type Contact struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}
