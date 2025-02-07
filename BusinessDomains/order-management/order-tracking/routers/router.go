package routers

import (
	"github.com/gin-gonic/gin"
	"order-tracking/handlers"
)

func SetupRouter(router *gin.Engine) {
	tracking := router.Group("/tracking")
	{
		tracking.GET("/:id", handlers.TrackOrder)
	}
}
