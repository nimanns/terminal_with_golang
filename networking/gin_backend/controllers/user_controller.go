package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "gin_backend/models"
    "gin_backend/utils"
)

func Register(c *gin.Context) {
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.Password = utils.HashPassword(input.Password)

    if err := models.CreateUser(input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func Login(c *gin.Context) {
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := models.GetUserByEmail(input.Email)
    if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
