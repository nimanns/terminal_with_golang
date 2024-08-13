package main

import (
	"log"
	"net/http"

	"delicake/models"
	"delicake/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("delicake.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database successfully")

	err = db.AutoMigrate(&models.Cake{}, &models.Order{})
	if err != nil {
		log.Fatal("Failed to perform migrations:", err)
	}
	log.Println("Migrations completed successfully")

	cakeRepo := &repository.CakeRepository{DB: db}
	orderRepo := &repository.OrderRepository{DB: db}

	r := gin.Default()

	r.GET("/cakes", func(c *gin.Context) {
		cakes, err := cakeRepo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cakes)
	})

	r.GET("/orders", func(c *gin.Context) {
		orders, err := orderRepo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	})

	log.Fatal(r.Run(":8080"))
}
