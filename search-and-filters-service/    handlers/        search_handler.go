// handlers/search_handler.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "search-and-filters-service/services"
)

func SearchModels(c *gin.Context) {
    query := c.Query("query")
    results, err := services.SearchModels(query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, results)
}
