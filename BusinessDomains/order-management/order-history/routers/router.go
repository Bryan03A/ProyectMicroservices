package routers

import (
	"github.com/gin-gonic/gin"
	"order-history/handlers"
)

func SetupRouter(router *gin.Engine) {
	history := router.Group("/history")
	{
		history.GET("/:user_id", handlers.GetOrderHistory)
	}
}
