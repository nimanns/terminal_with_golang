package routes

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "gin_backend/controllers"
    "gin_backend/utils"
)

func InitializeRoutes(router *gin.Engine, db *gorm.DB) {
    router.POST("/register", controllers.Register)
    router.POST("/login", controllers.Login)

    authorized := router.Group("/")
    authorized.Use(utils.AuthMiddleware())
    {
        authorized.POST("/upload", controllers.UploadPhoto)
    }
}
