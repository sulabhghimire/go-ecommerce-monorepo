package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
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

	pvtRoutes.Get("/cart", handler.addToCart)
	pvtRoutes.Post("/cart", handler.getCart)
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

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user profile fetched successfully",
		"data":    user,
	})

}

func (h UserHandler) createProfile(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "register",
	})

}

func (h UserHandler) addToCart(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "register",
	})

}

func (h UserHandler) getCart(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "register",
	})

}

func (h UserHandler) createOrder(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "register",
	})

}

func (h UserHandler) getOrders(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "register",
	})

}

func (h UserHandler) getOrder(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "register",
	})

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
