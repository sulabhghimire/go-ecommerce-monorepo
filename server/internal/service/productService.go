package service

import (
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
)

type ProductService struct {
	Repo repository.ProductRepository
	Auth	helper.Auth
	Config config.AppConfig
}

func (s ProductService) CreateProduct (input dto.CreateProductRequest, user domain.User) (*domain.Product,error) {
	return s.Repo.CreateProduct(&domain.Product{
		Name: input.Name,
		Description: input.Description,
		Price: input.Price,
		CategoryId: input.CategoryId,
		ImageUrl: input.ImageUrl,
		UserId: user.ID,
		Stock: uint(input.Stock),
	})
}

func (s ProductService) EditProduct (id uint,input dto.CreateProductRequest, user domain.User) (*domain.Product,error) {

	currProd, err := s.Repo.GetProductById(id);
	if err != nil {
		return nil, err
	}

	if currProd.UserId != user.ID {
		return  nil, helper.NOT_AUTHORIZED_ERROR
	}

	if len(input.Name) > 0 {
		currProd.Name = input.Name
	}
	if len(input.Description) > 0 {
		currProd.Description = input.Description
	}
	if len(input.ImageUrl) > 0 {
		currProd.ImageUrl = input.ImageUrl
	}
	if input.Price > 0 {
		currProd.Price = input.Price
	}
	if input.CategoryId > 0 {
		currProd.CategoryId = input.CategoryId
	}
	if input.Stock > 0 {
		currProd.Stock = uint(input.Stock)
	}

	return s.Repo.EditProduct(currProd)

}

func (s ProductService) DeleteProduct(id uint, user domain.User) error {

	prod, err := s.Repo.GetProductById(id)
	if err != nil {
		return  err
	}

	if prod.UserId != user.ID {
		return  helper.NOT_AUTHORIZED_ERROR
	}

	return s.Repo.DeleteProduct(id);
	
} 

func (s ProductService) GetProducts()([]*domain.Product, error) {

	return s.Repo.GetProducts();

}

func (s ProductService) GetProductById(id uint) (*domain.Product, error) {

	return s.Repo.GetProductById(id);

}

func (s ProductService) GetSellerProducts(sellerId uint) ([]*domain.Product, error) {
	return s.Repo.FindSellerProducts(sellerId)
}

func (s ProductService) UpdateProductStock(id uint, stock int, user domain.User) (*domain.Product, error) {
	prod, err := s.Repo.GetProductById(id);
	if err != nil {
		return nil, err
	}

	if prod.UserId != user.ID {
		return nil, helper.NOT_AUTHORIZED_ERROR
	}

	prod.Stock = uint(stock)
	updatedProd, err := s.Repo.EditProduct(prod)
	if err != nil {
		return nil, err
	}

	return updatedProd, nil
}