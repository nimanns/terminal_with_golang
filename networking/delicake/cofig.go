package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func NewConfig() *Config {
	db, err := gorm.Open(sqlite.Open("delicake.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database successfully")

	return &Config{
		DB: db,
	}
}
