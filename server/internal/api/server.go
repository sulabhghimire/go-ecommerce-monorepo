package api

import (
	"ecommerce/config"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/domain"
	"ecommerce/internal/helper"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
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
	err = db.AutoMigrate(&domain.User{}, &domain.BankAccount{}, &domain.Category{}, &domain.Product{})
	if err != nil {
		log.Fatalln("error on  migration %v", err.Error())
	}

	log.Println("migration done successfully")

	fmt.Println(config.AppSecret)
	auth := helper.SetUpAuth(config.AppSecret)

	restHandler := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
	}

	setUpRoutes(restHandler)

	app.Listen(config.ServerPort)
}

func setUpRoutes(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)
	handlers.SetUpCatalogRoutes(rh)
}
