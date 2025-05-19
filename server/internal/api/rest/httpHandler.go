package rest

import (
	"ecommerce/config"
	"ecommerce/internal/helper"
	"ecommerce/pkg/payment"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config config.AppConfig
	Pc     payment.PaymentClient
}
