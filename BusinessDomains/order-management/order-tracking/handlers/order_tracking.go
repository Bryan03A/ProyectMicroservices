package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"order-tracking/models"
	"order-tracking/services"
)

func TrackOrder(c *gin.Context) {
	id := c.Param("id")
	status, err := services.TrackOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}
