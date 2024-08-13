package models

import "gorm.io/gorm"

type Cake struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
}

