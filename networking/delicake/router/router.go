package router

import (
	"delicake/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cakeHandler *handlers.CakeHandler, orderHandler *handlers.OrderHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/cakes", cakeHandler.CreateCake)
	r.GET("/cakes", cakeHandler.GetAllCakes)

	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders", orderHandler.GetAllOrders)

	return r
}
