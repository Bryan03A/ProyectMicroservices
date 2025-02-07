package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"order-history/models"
	"order-history/services"
)

func GetOrderHistory(c *gin.Context) {
	userID := c.Param("user_id")
	history, err := services.GetOrderHistory(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
