package repository

import (
	"ecommerce/internal/domain"
	"errors"
	"log"

	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateCategory(e *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryById(id uint) (*domain.Category, error)
	EditCategory(e *domain.Category) (*domain.Category, error)
	DeleteCategory(id uint) error
}

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {

	return &catalogRepository{
		db: db,
	}

}

func (r catalogRepository) CreateCategory(e *domain.Category) error {

	err := r.db.Create(&e).Error
	if err != nil {
		log.Printf("Create category failed %v.\n", err)
		return errors.New("create category failed")
	}

	return nil
}

func (r catalogRepository) FindCategories() ([]*domain.Category, error) {

	var categories []*domain.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r catalogRepository) FindCategoryById(id uint) (*domain.Category, error) {

	var category *domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		log.Printf("Find category failed %v.\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrorCategoryNotFound
		}
		return nil, errors.New("failed to find category of given id")
	}

	return category, nil
}

func (r catalogRepository) EditCategory(e *domain.Category) (*domain.Category, error) {

	err := r.db.Save(&e).Error
	if err != nil {
		log.Printf("Find to update category %v.\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrorCategoryNotFound
		}
		return nil, errors.New("failed to update category of given id")
	}

	return e, nil
}

func (r catalogRepository) DeleteCategory(id uint) error {

	err := r.db.Delete(&domain.Category{}, id).Error
	if err != nil {
		log.Printf("Find to delete category %v.\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrorCategoryNotFound
		}
		return errors.New("failed to delete category of given id")
	}

	return nil
}
