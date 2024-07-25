package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "gin_backend/models"
)

func UploadPhoto(c *gin.Context) {
    file, _ := c.FormFile("file")

    photo := models.Photo{
        Filename: file.Filename,
        UserID:   c.MustGet("userID").(uint),
    }

    if err := models.CreatePhoto(photo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.SaveUploadedFile(file, "./uploads/"+file.Filename)
    c.JSON(http.StatusOK, gin.H{"message": "upload successful"})
}
