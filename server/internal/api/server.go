package api

import (
	"ecommerce/config"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/domain"
	"ecommerce/internal/helper"
	"ecommerce/pkg/payment"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {

	app := fiber.New()

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection error %v\n", err)
	}
	log.Println("Database connected successfully")

	// run migration
	err = db.AutoMigrate(
		&domain.User{},
		&domain.BankAccount{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Address{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{})
	if err != nil {
		log.Fatalf("error on  migration %v", err.Error())
	}
	log.Println("migration done successfully")

	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000/",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})

	app.Use(c)

	auth := helper.SetUpAuth(config.AppSecret)

	paymentClient := payment.NewPaymentClient(config.StripeConfig.StripeSecretKey, config.StripeConfig.SuccessUrl, config.StripeConfig.CancelUrl)

	restHandler := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
		Pc:     paymentClient,
	}

	setUpRoutes(restHandler)

	app.Listen(config.ServerPort)
}

func setUpRoutes(rh *rest.RestHandler) {
	handlers.SetUpCatalogRoutes(rh)
	handlers.SetupUserRoutes(rh)
	handlers.SetupTransactionRoutes(rh)
}
