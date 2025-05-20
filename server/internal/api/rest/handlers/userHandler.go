package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	// svc UserService
	svc service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {

	app := rh.App

	svc := service.UserService{
		Repo:   repository.NewUserRepository(rh.DB),
		PRepo:  repository.NewProductRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}
	handler := UserHandler{
		svc: svc,
	}

	// Public endpoints
	pubRoutes := app.Group("/users")

	pubRoutes.Post("/register", handler.register)
	pubRoutes.Post("/login", handler.login)

	// Private endpoints

	pvtRoutes := pubRoutes.Group("/", rh.Auth.Authorize)

	pvtRoutes.Get("/verify", handler.getVerificationCode)
	pvtRoutes.Post("/verify", handler.verify)

	pvtRoutes.Get("/profile", handler.getProfile)
	pvtRoutes.Post("/profile", handler.createProfile)
	pvtRoutes.Patch("/profile", handler.updateProfile)

	pvtRoutes.Get("/cart", handler.getCart)
	pvtRoutes.Post("/cart", handler.addToCart)

	pvtRoutes.Post("/order", handler.createOrder)
	pvtRoutes.Get("/order", handler.getOrders)
	pvtRoutes.Get("/order/:id", handler.getOrder)

	pvtRoutes.Post("/become-seller", handler.becomeSeller)

}

func (h UserHandler) register(ctx *fiber.Ctx) error {

	user := dto.UserSignUp{}

	err := ctx.BodyParser(&user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid input",
		})
	}

	token, err := h.svc.SignUp(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error on signup",
			"reason":  err,
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user signup sucessfull",
		"token":   token,
	})

}

func (h UserHandler) login(ctx *fiber.Ctx) error {

	loginInInput := dto.UserLogin{}

	err := ctx.BodyParser(&loginInInput)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid input",
		})
	}

	token, err := h.svc.Login(loginInInput.Email, loginInInput.Password)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide correct user email and password",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user login sucessfull",
		"token":   token,
	})

}

func (h UserHandler) getVerificationCode(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	err := h.svc.GetVerificationCode(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "verification code generated",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "verification code has been sent your contact number",
	})

}

func (h UserHandler) verify(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	payload := dto.VerificationCodeInput{}
	err := ctx.BodyParser(&payload)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid input",
		})
	}

	err = h.svc.VerifyCode(user.ID, payload.Code)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "error verifying the code",
			"reason":  err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user verified successfully",
	})

}

func (h UserHandler) getProfile(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrorUserNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Profile fetched sucessfully", profile)

}

func (h UserHandler) updateProfile(ctx *fiber.Ctx) error {

	payload := dto.ProfileInput{}
	err := ctx.BodyParser(&payload)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid input",
		})
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	_, err = h.svc.UpdateProfile(user.ID, payload)
	if err != nil {
		if errors.Is(err, domain.ErrorUserNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Profile updated sucessfully", nil)

}

func (h UserHandler) createProfile(ctx *fiber.Ctx) error {

	payload := dto.ProfileInput{}
	if err := ctx.BodyParser(&payload); err != nil {
		return rest.BadRequest(ctx, "please provide valid input")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	err := h.svc.CreateProfile(user.ID, payload)
	if err != nil {
		if errors.Is(err, domain.ErrorUserNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Profile created sucessfully", nil)

}

func (h UserHandler) addToCart(ctx *fiber.Ctx) error {

	paylaod := dto.CreateCartRequest{}
	if err := ctx.BodyParser(&paylaod); err != nil {
		return rest.BadRequest(ctx, "please provide valid payload")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	cartItems, err := h.svc.CreateCart(paylaod, user)
	if err != nil {
		if errors.Is(err, domain.ProductNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Item added to cart successfully", cartItems)

}

func (h UserHandler) getCart(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	cart, _, err := h.svc.FindCart(user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Cart fetched sucessfully", cart)

}

func (h UserHandler) createOrder(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	orderRef, err := h.svc.CreateOrder(user)
	if err != nil {
		if errors.Is(err, domain.ErrorCartItemNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusCreated, "Order created sucessfully", map[string]interface{}{"order_ref": orderRef})
}

func (h UserHandler) getOrders(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	orders, err := h.svc.GetOrders(user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Orders fetched sucessfully", orders)

}

func (h UserHandler) getOrder(ctx *fiber.Ctx) error {

	orderId := ctx.Params("id")
	if orderId == "" {
		return rest.BadRequest(ctx, "please provide a valid order id")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	order, err := h.svc.GetOrderById(orderId, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrorOrderNotFound) {
			return rest.NotFoundError(ctx, err)
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, http.StatusOK, "Order fetched sucessfully", order)

}

func (h UserHandler) becomeSeller(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	payload := dto.SellerInput{}
	err := ctx.BodyParser(&payload)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid input",
		})
	}

	token, err := h.svc.BecomeSeller(user.ID, payload)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "error occurred while updating to seller status",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "updated to seller successfully",
		"token":   token,
	})

}
