// handlers/index_handler.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "search-and-filters-service/services"
)

func IndexModel(c *gin.Context) {
    var modelData map[string]interface{}
    if err := c.ShouldBindJSON(&modelData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    modelID := modelData["model_id"].(string)
    err := services.IndexModel(modelID, modelData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"result": "Modelo indexado"})
}
