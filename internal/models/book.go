package models

import "time"

type Book struct {
	ID        int64     `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Title     string    `json:"title" binding:"required"`
	Author    string    `json:"author" binding:"required"`
	ISBN      string    `json:"isbn"`
	ImageURL  string    `json:"image_url"`
	Price     float64   `json:"price" binding:"required,min=0"`
	Stock     int       `json:"stock" binding:"required,min=0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
