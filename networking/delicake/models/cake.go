package models

import (
	"gorm.io/gorm"
)

type Cake struct {
	gorm.Model
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description" binding:"max=500"`
	Price       float64 `json:"price" binding:"required,min=0"`
}
