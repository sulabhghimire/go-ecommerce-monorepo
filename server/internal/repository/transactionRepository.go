package repository

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindInitialPayment(u uint) (*domain.Payment, error)
	FindOrders(uID uint) ([]domain.Order, error)
	FindOrderById(orderId, uID uint) (dto.SellerOrderDetails, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// FindPayment implements TransactionRepository.
func (r *transactionRepository) FindInitialPayment(u uint) (*domain.Payment, error) {
	var payment *domain.Payment
	err := r.db.First(&payment, "user_id=? AND status=initial", u).Order("created_at desc").Error
	if err != nil {
		fmt.Printf("data base error cccured %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrorUserInitialPaymentNotFound
		}
		return nil, errors.New("some error occured")
	}
	return payment, nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreatePayment(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *transactionRepository) FindOrders(uID uint) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *transactionRepository) FindOrderById(orderId, uID uint) (dto.SellerOrderDetails, error) {
	return dto.SellerOrderDetails{}, nil
}
