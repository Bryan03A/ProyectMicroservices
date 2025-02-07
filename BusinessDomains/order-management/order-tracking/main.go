package main

import (
	"order-tracking/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.SetupRouter(r)
	r.Run(":8082")
}
