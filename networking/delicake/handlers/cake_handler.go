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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	cakes, total, err := h.service.GetAllCakes(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cakes,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// handlers/order_handler.go
// (Similar structure to cake_handler.go, with CreateOrder and GetAllOrders methods)
