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
