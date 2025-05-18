package repository

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uID uint) ([]domain.Order, error)
	FindOrderById(orderId, uID uint) (dto.SellerOrderDetails, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreatePayment(payment *domain.Payment) error {
	return nil
}

func (r *transactionRepository) FindOrders(uID uint) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *transactionRepository) FindOrderById(orderId, uID uint) (dto.SellerOrderDetails, error) {
	return dto.SellerOrderDetails{}, nil
}
