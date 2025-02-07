package main

import (
	"order-history/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.SetupRouter(r)
	r.Run(":8081")
}
