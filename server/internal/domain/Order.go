package domain

import (
	"errors"
	"time"
)

var (
	ErrorOrderNotFound = errors.New("order not found")
)

type Order struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserId        uint        `json:"user_id"`
	Status        string      `json:"status"`
	Amount        float64     `json:"amount"`
	TransactionId string      `json:"transaction_id"`
	OrderRef      string      `json:"order_ref" gorm:"index;unique;not null"`
	PaymentId     string      `json:"payment_id"`
	Itens         []OrderItem `json:"items"`
	CreatedAt     time.Time   `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"default:current_timestamp"`
}
