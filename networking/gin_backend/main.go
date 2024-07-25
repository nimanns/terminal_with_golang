package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "gin_backend/config"
    "gin_backend/routes"
    "gin_backend/models"
)

func main() {
    db := initDatabase()
    defer db.Close()

    router := gin.Default()
    routes.InitializeRoutes(router, db)

    router.Run(":8080")
}

func initDatabase() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&models.User{}, &models.Photo{})
    return db
}
