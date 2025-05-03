package main

import (
	"ecommerce/config"
	"ecommerce/internal/api"
	"log"
)

func main() {

	cfg, err := config.SetUpEnv()
	if err != nil {
		log.Fatalf("Config file is not loaded properly %v\n", err)
	}

	api.StartServer(cfg)

}
