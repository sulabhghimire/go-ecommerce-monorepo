package domain

import (
	"errors"
	"time"
)

var (
	ErrorProductNotFound   = errors.New("product of given id not found")
	ErrorStockNotAvailable = errors.New("stock not available")
)

type Product struct {
	ID          uint      `json:"id" gorm:"PrimaryKey"`
	Name        string    `json:"name" gorm:"index;"`
	Description string    `json:"description"`
	CategoryId  uint      `json:"category_id"`
	ImageUrl    string    `json:"image_url"`
	Price       float64   `json:"price"`
	UserId      uint      `json:"user_id"`
	Stock       uint      `json:"stock"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
