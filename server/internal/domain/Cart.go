package domain

import (
	"errors"
	"time"
)

var (
	ErrorUserProductCartNotFound = errors.New("cart of given user and product not found")
	ErrorCartItemNotFound        = errors.New("cart item not found")
)

type Cart struct {
	ID        uint      `gorm:"PrimaryKey" json:"id"`
	UserId    uint      `json:"user_id"`
	ProductId uint      `json:"product_id"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
	SellerId  uint      `json:"seller_id"`
	Price     float64   `json:"price"`
	Qty       uint      `json:"qty"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
