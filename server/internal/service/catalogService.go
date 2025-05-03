package service

import (
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryDTO) error {

	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	})

	return err
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {

	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s CatalogService) GetCategory(id uint) (*domain.Category, error) {

	cat, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (s CatalogService) DeleteCategory(id uint) error {

	return s.Repo.DeleteCategory(id)
}

func (s CatalogService) EditCategory(id uint, input dto.CreateCategoryDTO) (*domain.Category, error) {

	exCat, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, err
	}

	if len(input.Name) > 0 {
		exCat.Name = input.Name
	}
	if input.ParentId > 0 {
		exCat.ParentId = input.ParentId
	}
	if len(input.ImageUrl) > 0 {
		exCat.ImageUrl = input.ImageUrl
	}
	if input.DisplayOrder > 0 {
		exCat.DisplayOrder = input.DisplayOrder
	}

	updatedCat, err := s.Repo.EditCategory(exCat)
	if err != nil {
		return nil, err
	}

	return updatedCat, nil
}
