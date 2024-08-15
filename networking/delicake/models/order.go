package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerName string  `json:"customer_name" binding:"required,min=2,max=100"`
	CakeID       uint    `json:"cake_id" binding:"required"`
	Cake         Cake    `json:"cake"`
	Quantity     int     `json:"quantity" binding:"required,min=1"`
	TotalPrice   float64 `json:"total_price"`
}
