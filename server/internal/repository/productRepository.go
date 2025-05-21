package repository

import (
	"ecommerce/internal/domain"
	"errors"
	"log"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(*domain.Product) (*domain.Product, error)
	GetProducts() ([]*domain.Product, error)
	GetProductById(id uint) (*domain.Product, error)
	EditProduct(*domain.Product) (*domain.Product, error)
	DeleteProduct(id uint) error
	FindSellerProducts(sellerId uint) ([]*domain.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

// FindSellerProducts implements ProductRepository.
func (p *productRepository) FindSellerProducts(sellerId uint) ([]*domain.Product, error) {

	var products []*domain.Product
	err := p.db.Where("user_id=?", sellerId).Find(&products).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return nil, errors.New("error fetching products")
	}

	return products, nil

}

// CreateProduct implements ProductRepository.
func (p productRepository) CreateProduct(e *domain.Product) (*domain.Product, error) {

	err := p.db.Model(&domain.Product{}).Create(e).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return nil, errors.New("error creating product")
	}

	return e, nil
}

// DeleteProduct implements ProductRepository.
func (p productRepository) DeleteProduct(id uint) error {

	err := p.db.Delete(&domain.Product{}, id).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrorProductNotFound
		}
		return errors.New("error deleting product")
	}

	return nil

}

// EditProduct implements ProductRepository.
func (p productRepository) EditProduct(e *domain.Product) (*domain.Product, error) {

	err := p.db.Save(&e).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrorProductNotFound
		}
		return nil, errors.New("error deleting product")
	}

	return e, nil
}

// GetProductById implements ProductRepository.
func (p productRepository) GetProductById(id uint) (*domain.Product, error) {

	var product *domain.Product
	err := p.db.First(&product, id).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrorProductNotFound
		}
		return nil, errors.New("error fetching product of given id")
	}

	return product, nil

}

// GetProducts implements ProductRepository.
func (p productRepository) GetProducts() ([]*domain.Product, error) {

	var product []*domain.Product
	err := p.db.Find(&product).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return nil, errors.New("error fetching product of given id")
	}

	return product, nil
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}
