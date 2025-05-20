package service

import (
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"

	"github.com/stripe/stripe-go/v78"
)

type TransactionService struct {
	Auth   helper.Auth
	Config config.AppConfig
	Repo   repository.TransactionRepository
}

func (s TransactionService) GetActivePayments(uId uint) (*domain.Payment, error) {
	return s.Repo.FindInitialPayment(uId)
}

func (s TransactionService) StoreCreatedPayment(uId uint, ps *stripe.CheckoutSession, amount float64, orderId string) error {

	payment := domain.Payment{
		UserId:     uId,
		Amount:     amount,
		Status:     domain.PaymentStatusInitial,
		PaymentUrl: ps.URL,
		PaymentId:  ps.ID,
		OrderId:    orderId,
	}

	return s.Repo.CreatePayment(&payment)
}

func (s TransactionService) GetOrders(u uint) ([]domain.Order, error) {
	orders, err := s.Repo.FindOrders(u)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s TransactionService) GetOrderById(orderId, sellerId uint) (dto.SellerOrderDetails, error) {
	order, err := s.Repo.FindOrderById(orderId, sellerId)
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}
	return order, nil
}
