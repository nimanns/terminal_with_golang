package main

import (
	"log"

	"delicake/config"
	"delicake/handlers"
	"delicake/models"
	"delicake/repository"
	"delicake/router"
	"delicake/service"
)

func main() {
	cfg := config.NewConfig()

	err := cfg.DB.AutoMigrate(&models.Cake{}, &models.Order{})
	if err != nil {
		log.Fatal("Failed to perform migrations:", err)
	}
	log.Println("Migrations completed successfully")

	cakeRepo := repository.NewCakeRepository(cfg.DB)
	orderRepo := repository.NewOrderRepository(cfg.DB)

	cakeService := service.NewCakeService(cakeRepo)
	orderService := service.NewOrderService(orderRepo)

	cakeHandler := handlers.NewCakeHandler(cakeService)
	orderHandler := handlers.NewOrderHandler(orderService)

	r := router.SetupRouter(cakeHandler, orderHandler)

	log.Fatal(r.Run(":8080"))
}
