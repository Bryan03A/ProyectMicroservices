// main.go
package main

import (
    "github.com/gin-gonic/gin"
    _ "search-and-filters-service/config"
    "search-and-filters-service/handlers"
)

func main() {
    r := gin.Default()

    r.POST("/index", handlers.IndexModel)
    r.GET("/search", handlers.SearchModels)

    r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
