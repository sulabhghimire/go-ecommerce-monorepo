package domain

import "time"

type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserId        uint      `json:"user_id"`
	CaptureMethod string    `json:"capture_method"`
	Amount        float64   `json:"amount"`
	TransactionId string    `json:"transaction_id"`
	CustomerId    string    `json:"customer_id"` // stripe id
	PaymentId     string    `json:"payment_id"`  // paymnent id
	Status        string    `json:"status"`      // initial, success, failed
	Response      string    `json:"response"`    // response from payment gateway
	CreatedAt     time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
