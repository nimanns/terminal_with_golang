package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CustomerName string
	CakeID       uint
	Cake         Cake
	Quantity     int
}
