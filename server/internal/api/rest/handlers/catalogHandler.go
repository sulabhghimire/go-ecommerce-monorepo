package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CatalogHandler struct {
	svc service.CatalogService
}

func SetUpCatalogRoutes(rh *rest.RestHandler) {

	app := rh.App

	catalogRepo := repository.NewCatalogRepository(rh.DB)

	svc := service.CatalogService{
		Auth:   rh.Auth,
		Repo:   catalogRepo,
		Config: rh.Config,
	}

	handler := CatalogHandler{
		svc: svc,
	}

	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)
	app.Get("/categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategoryById)

	sellerRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)
	// categories
	sellerRoutes.Post("/categories", handler.CreateCategories)
	sellerRoutes.Patch("/categories/:id", handler.EditCategories)
	sellerRoutes.Delete("/categories/:id", handler.DeleteCategories)
	// products
	sellerRoutes.Post("/products", handler.CreateProducts)
	sellerRoutes.Get("/products", handler.GetProducts)
	sellerRoutes.Get("/products/:id", handler.GetProduct)
	sellerRoutes.Patch("/products/:id", handler.UpdateStock)
	sellerRoutes.Delete("/products/:id", handler.DeleteProducts)
	sellerRoutes.Put("/products/:id", handler.EditProducts)

}

func (h CatalogHandler) GetCategories(ctx *fiber.Ctx) error {

	categories, err := h.svc.GetCategories()
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Categories fetched successfully", categories)

}

func (h CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {

	catId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequest(ctx, "please provide a valid category id")
	}

	cat, err := h.svc.GetCategory(uint(catId))
	if err != nil {
		if errors.Is(err, domain.CategoryNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Category fetched successfully", cat)

}

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {

	payload := dto.CreateCategoryDTO{}
	if err := ctx.BodyParser(&payload); err != nil {
		log.Printf("valid input %v", err)
		return rest.BadRequest(ctx, "please provide proper valid input")
	}

	err := h.svc.CreateCategory(payload)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Category created successfully", nil)
}

func (h CatalogHandler) EditCategories(ctx *fiber.Ctx) error {

	payload := dto.CreateCategoryDTO{}
	if err := ctx.BodyParser(&payload); err != nil {
		log.Printf("valid input %v", err)
		return rest.BadRequest(ctx, "please provide proper valid input")
	}

	catId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || catId < 0 {
		return rest.BadRequest(ctx, "please provide a valid category id")
	}

	cat, err := h.svc.EditCategory(uint(catId), payload)
	if err != nil {
		if errors.Is(err, domain.CategoryNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Category edited successfully", cat)
}

func (h CatalogHandler) DeleteCategories(ctx *fiber.Ctx) error {

	catId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || catId < 0 {
		return rest.BadRequest(ctx, "please provide a valid category id")
	}

	err = h.svc.DeleteCategory(uint(catId))
	if err != nil {
		if errors.Is(err, domain.CategoryNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Category deleted successfully", nil)
}

func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, http.StatusOK, "Product created successfully", nil)
}

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, http.StatusOK, "Products fetched successfully", nil)
}

func (h CatalogHandler) GetProduct(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, http.StatusOK, "Product fetched successfully", nil)
}

func (h CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, http.StatusOK, "Product Stock updated successfully", nil)
}

func (h CatalogHandler) DeleteProducts(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, http.StatusOK, "Product deleted successfully", nil)
}

func (h CatalogHandler) EditProducts(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, http.StatusOK, "Product edited successfully", nil)
}
