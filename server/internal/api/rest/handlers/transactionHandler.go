package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/payment"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	svc           service.TransactionService
	paymentClient payment.PaymentClient
}

func SetupTransactionRoutes(rh *rest.RestHandler) {

	app := rh.App

	transactionRepo := repository.NewTransactionRepository(rh.DB)
	svc := service.TransactionService{
		Auth:   rh.Auth,
		Config: rh.Config,
		Repo:   transactionRepo,
	}

	handler := TransactionHandler{
		svc:           svc,
		paymentClient: rh.Pc,
	}

	pvtRoutes := app.Group("/transactions", rh.Auth.Authorize)
	pvtRoutes.Get("/payment", handler.MakePayment)

	sellerRoutes := app.Group("/transactions/seller", rh.Auth.AuthorizeSeller)
	sellerRoutes.Get("/orders", handler.GetOrders)
	sellerRoutes.Get("/orders/:id", handler.GetOrderById)
}

func (h *TransactionHandler) MakePayment(c *fiber.Ctx) error {

	return rest.SuccessResponse(c, fiber.StatusOK, "Payment successful", nil)
}

func (h *TransactionHandler) GetOrders(c *fiber.Ctx) error {

	return rest.SuccessResponse(c, fiber.StatusOK, "Orders fetched successfully", nil)
}

func (h *TransactionHandler) GetOrderById(c *fiber.Ctx) error {
	return rest.SuccessResponse(c, fiber.StatusOK, "Order fetched successfully", nil)
}
