package main

import (
	"log"

	"delicake/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("delicake.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database successfully")

	// migrations
	err = db.AutoMigrate(&models.Cake{}, &models.Order{})
	if err != nil {
		log.Fatal("Failed to perform migrations:", err)
	}
	log.Println("Migrations completed successfully")
}
