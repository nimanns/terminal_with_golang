package handlers

import (
	"delicake/models"
	"delicake/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CakeHandler struct {
	service *service.CakeService
}

func NewCakeHandler(service *service.CakeService) *CakeHandler {
	return &CakeHandler{service: service}
}

func (h *CakeHandler) CreateCake(c *gin.Context) {
	var cake models.Cake
	if err := c.ShouldBindJSON(&cake); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCake(&cake); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cake)
}

func (h *CakeHandler) GetAllCakes(c *gin.Context) {
	cakes, err := h.service.GetAllCakes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cakes)
}

// handlers/order_handler.go
// (Similar structure to cake_handler.go, with CreateOrder and GetAllOrders methods)
