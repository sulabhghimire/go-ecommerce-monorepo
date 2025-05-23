package handlers

import (
	"ecommerce/config"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/payment"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	svc           service.TransactionService
	userSvc       service.UserService
	paymentClient payment.PaymentClient
	cfg           config.AppConfig
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
		cfg:           rh.Config,
	}

	pvtRoutes := app.Group("/transactions", rh.Auth.Authorize)
	pvtRoutes.Get("/payment", handler.MakePayment)
	pvtRoutes.Get("/payment/verify", handler.VerifyPayment)

	sellerRoutes := app.Group("/transactions/seller", rh.Auth.AuthorizeSeller)
	sellerRoutes.Get("/orders", handler.GetOrders)
	sellerRoutes.Get("/orders/:id", handler.GetOrderById)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	pubKey := h.cfg.StripeConfig.PublishableKey

	activePayment, err := h.svc.GetActivePayments(user.ID)
	if err != nil {
		fmt.Println(errors.Is(err, domain.ErrorUserInitialPaymentNotFound))
		if !errors.Is(err, domain.ErrorUserInitialPaymentNotFound) {
			return rest.InternalError(ctx, err)
		}
	}

	if activePayment != nil {
		return rest.SuccessResponse(ctx, http.StatusOK, "payment session created", map[string]interface{}{
			"pubKey": pubKey,
			"secret": activePayment.ClientSecret,
		})
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

	paymentResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	err = h.svc.StoreCreatedPayment(dto.CreatePaymentRequest{
		UserId:       user.ID,
		OrderId:      orderId,
		Amount:       amount,
		ClientSecret: paymentResult.ClientSecret,
		PaymentId:    paymentResult.ID,
	})
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, fiber.StatusOK, "Payment session created", map[string]interface{}{
		"pubkey": pubKey,
		"secret": paymentResult.ClientSecret,
	})
}

func (h *TransactionHandler) VerifyPayment(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	activePayment, err := h.svc.GetActivePayments(user.ID)
	if err != nil || activePayment == nil {
		if errors.Is(err, domain.ErrorUserInitialPaymentNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	paymentRes, err := h.paymentClient.GetPaymentStatus(activePayment.PaymentId)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	paymentJson, _ := json.Marshal(paymentRes)
	paymentLogs := string(paymentJson)
	paymentStatus := "failed"
	msg := "Payment failed"

	if paymentRes.Status == "succeeded" {
		// Create order here
		paymentStatus = "success"
		msg = "Payment verified sucessfully"
		err = h.userSvc.CreateOrder(user.ID, activePayment.OrderId, activePayment.PaymentId, activePayment.Amount)
		if err != nil {
			return rest.InternalError(ctx, err)
		}
	}

	h.svc.UpdatePayment(user.ID, paymentStatus, paymentLogs)

	return rest.SuccessResponse(ctx, fiber.StatusOK, msg, nil)

}

func (h *TransactionHandler) GetOrders(c *fiber.Ctx) error {

	return rest.SuccessResponse(c, fiber.StatusOK, "Orders fetched successfully", nil)
}

func (h *TransactionHandler) GetOrderById(c *fiber.Ctx) error {
	return rest.SuccessResponse(c, fiber.StatusOK, "Order fetched successfully", nil)
}
