package service

import (
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
)

type TransactionService struct {
	Auth   helper.Auth
	Config config.AppConfig
	Repo   repository.TransactionRepository
}

func (s TransactionService) GetActivePayments(uId uint) (*domain.Payment, error) {
	return s.Repo.FindInitialPayment(uId)
}

func (s TransactionService) StoreCreatedPayment(payload dto.CreatePaymentRequest) error {

	payment := domain.Payment{
		UserId:       payload.UserId,
		Amount:       payload.Amount,
		Status:       domain.PaymentStatusInitial,
		ClientSecret: payload.ClientSecret,
		PaymentId:    payload.PaymentId,
		OrderId:      payload.OrderId,
	}

	return s.Repo.CreatePayment(&payment)
}

func (s TransactionService) UpdatePayment(userId uint, status string, paymentLog string) error {
	p, err := s.GetActivePayments(userId)
	if err != nil {
		return err
	}

	p.Status = domain.PaymentStatus(status)
	p.Response = paymentLog

	return s.Repo.UpdatePayment(p)

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
