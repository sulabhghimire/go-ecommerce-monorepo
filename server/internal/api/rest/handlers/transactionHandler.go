package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/domain"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/payment"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	svc           service.TransactionService
	userSvc       service.UserService
	paymentClient payment.PaymentClient
}

func SetupTransactionRoutes(rh *rest.RestHandler) {

	app := rh.App

	transactionRepo := repository.NewTransactionRepository(rh.DB)
	userRepo := repository.NewUserRepository(rh.DB)
	productRepo := repository.NewProductRepository(rh.DB)

	userSvc := service.UserService{
		Repo:   userRepo,
		Auth:   rh.Auth,
		PRepo:  productRepo,
		Config: rh.Config,
	}

	svc := service.TransactionService{
		Auth:   rh.Auth,
		Config: rh.Config,
		Repo:   transactionRepo,
	}

	handler := TransactionHandler{
		svc:           svc,
		userSvc:       userSvc,
		paymentClient: rh.Pc,
	}

	pvtRoutes := app.Group("/transactions", rh.Auth.Authorize)
	pvtRoutes.Get("/payment", handler.MakePayment)

	sellerRoutes := app.Group("/transactions/seller", rh.Auth.AuthorizeSeller)
	sellerRoutes.Get("/orders", handler.GetOrders)
	sellerRoutes.Get("/orders/:id", handler.GetOrderById)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	activePayment, err := h.svc.GetActivePayments(user.ID)
	if err != nil {
		fmt.Println(errors.Is(err, domain.ErrorUserInitialPaymentNotFound))
		if !errors.Is(err, domain.ErrorUserInitialPaymentNotFound) {
			return rest.InternalError(ctx, err)
		}
	}

	if activePayment != nil {
		return rest.SuccessResponse(ctx, http.StatusOK, "payment session created", activePayment.PaymentUrl)
	}

	cartItems, amount, err := h.userSvc.FindCart(user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	if len(cartItems) == 0 {
		return rest.BadRequest(ctx, "You don't have any item to checkout")
	}

	orderId, err := helper.RandomString(8)
	if err != nil {
		return rest.InternalError(ctx, err)

	}

	sessionResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	err = h.svc.StoreCreatedPayment(user.ID, sessionResult, amount, orderId)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, fiber.StatusOK, "Payment session created", map[string]interface{}{
		"session_url": sessionResult.URL,
	})
}

func (h *TransactionHandler) GetOrders(c *fiber.Ctx) error {

	return rest.SuccessResponse(c, fiber.StatusOK, "Orders fetched successfully", nil)
}

func (h *TransactionHandler) GetOrderById(c *fiber.Ctx) error {
	return rest.SuccessResponse(c, fiber.StatusOK, "Order fetched successfully", nil)
}
