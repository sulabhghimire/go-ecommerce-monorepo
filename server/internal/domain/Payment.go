package domain

import (
	"errors"
	"time"
)

var (
	ErrorUserInitialPaymentNotFound = errors.New("inital payment not found of user")
)

type Payment struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	UserId        uint          `json:"user_id"`
	CaptureMethod string        `json:"capture_method"`
	Amount        float64       `json:"amount"`
	TransactionId string        `json:"transaction_id"`
	CustomerId    string        `json:"customer_id"` // stripe id
	PaymentId     string        `json:"payment_id"`  // paymnent id
	OrderId       string        `json:"order_id"`
	Status        PaymentStatus `json:"status" gorm:"default:initial"` // initial, success, failed, pending
	Response      string        `json:"response"`                      // response from payment gateway
	PaymentUrl    string        `json:"payment_url"`
	CreatedAt     time.Time     `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"default:current_timestamp"`
}

type PaymentStatus string

const (
	PaymentStatusInitial PaymentStatus = "inital"
	PaymensStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
)
