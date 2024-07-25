package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func CreateUser(user User) error {
    return db.Create(&user).Error
}

func GetUserByEmail(email string) (User, error) {
    var user User
    err := db.Where("email = ?", email).First(&user).Error
    return user, err
}
