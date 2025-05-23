package domain

import "time"

type OrderItem struct {
	ID        uint      `gorm:"PrimaryKey" json:"id"`
	ProductId uint      `json:"product_id"`
	UserId    uint      `json:"user_id"`
	OrderId   uint      `json:"order_id"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
	SellerId  uint      `json:"seller_id"`
	Price     float64   `json:"price"`
	Qty       uint      `json:"qty"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
