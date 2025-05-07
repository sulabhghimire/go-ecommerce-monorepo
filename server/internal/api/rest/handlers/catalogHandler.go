package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CatalogHandler struct {
	catalogSvc service.CatalogService
	prodSvc service.ProductService
}

func SetUpCatalogRoutes(rh *rest.RestHandler) {

	app := rh.App

	catalogRepo := repository.NewCatalogRepository(rh.DB)
	prodRepo := repository.NewProductRepository(rh.DB)

	catalogSvc := service.CatalogService{
		Auth:   rh.Auth,
		Repo:   catalogRepo,
		Config: rh.Config,
	}
	prodSvc := service.ProductService{
		Auth:   rh.Auth,
		Repo:   prodRepo,
		Config: rh.Config,
	}


	handler := CatalogHandler{
		catalogSvc: catalogSvc,
		prodSvc:prodSvc,
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
	sellerRoutes.Get("/products", handler.GetSellerProducts)
	sellerRoutes.Get("/products/:id", handler.GetProduct)
	sellerRoutes.Patch("/products/:id", handler.UpdateStock)
	sellerRoutes.Delete("/products/:id", handler.DeleteProducts)
	sellerRoutes.Put("/products/:id", handler.EditProducts)

}

func (h CatalogHandler) GetCategories(ctx *fiber.Ctx) error {

	categories, err := h.catalogSvc.GetCategories()
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

	cat, err := h.catalogSvc.GetCategory(uint(catId))
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

	err := h.catalogSvc.CreateCategory(payload)
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

	cat, err := h.catalogSvc.EditCategory(uint(catId), payload)
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

	err = h.catalogSvc.DeleteCategory(uint(catId))
	if err != nil {
		if errors.Is(err, domain.CategoryNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Category deleted successfully", nil)
}

func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {

	payload := dto.CreateProductRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		return rest.BadRequest(ctx, "please provide a valid request body")
	}

	user := h.prodSvc.Auth.GetCurrentUser(ctx)

	prod, err := h.prodSvc.CreateProduct(payload, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusCreated, "Product created successfully", prod)
}

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {

	prods, err := h.prodSvc.GetProducts()
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Products fetched successfully", prods)
}

func (h CatalogHandler) GetProduct(ctx *fiber.Ctx) error {

	prodId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || prodId < 0 {
		return rest.BadRequest(ctx, "please provide a valid product id")
	}

	prod, err := h.prodSvc.GetProductById(uint(prodId))
	if err != nil {
		if errors.Is(err, domain.ProductNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Product fetched successfully", prod)
}

func (h CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {

	prodId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || prodId < 0 {
		return rest.BadRequest(ctx, "please provide a valid product id")
	}

	paylaod := dto.UpdateStockRequest{}
	if err = ctx.BodyParser(&paylaod); err !=nil {
		return rest.BadRequest(ctx, "please provide a valid request body")
	}

	user := h.prodSvc.Auth.GetCurrentUser(ctx)

	prod, err := h.prodSvc.UpdateProductStock(uint(prodId), paylaod.Stock, user)
	if err != nil {
		if errors.Is(err, domain.ProductNotFound) {
			return rest.NotFoundError(ctx, err)
		}else if errors.Is(err, helper.NOT_AUTHORIZED_ERROR) {
			return rest.NotAuhtorizedError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Product Stock updated successfully", prod)
}

func (h CatalogHandler) DeleteProducts(ctx *fiber.Ctx) error {

	prodId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || prodId < 0 {
		return rest.BadRequest(ctx, "please provide a valid product id")
	}

	user := h.prodSvc.Auth.GetCurrentUser(ctx)

	err = h.prodSvc.DeleteProduct(uint(prodId),  user)
	if err != nil {
		if errors.Is(err, domain.ProductNotFound) {
			return rest.NotFoundError(ctx, err)
		}else if errors.Is(err, helper.NOT_AUTHORIZED_ERROR) {
			return rest.NotAuhtorizedError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Product deleted successfully", nil)
}

func (h CatalogHandler) EditProducts(ctx *fiber.Ctx) error {

	prodId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || prodId < 0 {
		return rest.BadRequest(ctx, "please provide a valid product id")
	}

	paylaod := dto.CreateProductRequest{}
	if err = ctx.BodyParser(&paylaod); err !=nil {
		return rest.BadRequest(ctx, "please provide a valid request body")
	}

	user := h.prodSvc.Auth.GetCurrentUser(ctx)

	prod, err := h.prodSvc.EditProduct(uint(prodId), paylaod, user)
	if err != nil {
		if errors.Is(err, domain.ProductNotFound) {
			return rest.NotFoundError(ctx, err)
		}else if errors.Is(err, helper.NOT_AUTHORIZED_ERROR) {
			return rest.NotAuhtorizedError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Product edited successfully", prod)
}

func (h CatalogHandler) GetSellerProducts(ctx *fiber.Ctx) error {

	user := h.prodSvc.Auth.GetCurrentUser(ctx)

	prods, err := h.prodSvc.GetSellerProducts(user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Products fetched successfully", prods)

}