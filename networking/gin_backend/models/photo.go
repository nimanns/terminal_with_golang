package models

import (
    "gorm.io/gorm"
)

type Photo struct {
    gorm.Model
    Filename string `json:"filename"`
    UserID   uint   `json:"user_id"`
}

func CreatePhoto(photo Photo) error {
    return db.Create(&photo).Error
}
