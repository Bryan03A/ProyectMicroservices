package handlers

import (
	"net/http"
	"order-management/order-service/src/models"
	"order-management/order-service/src/service"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderServiceInterface
}

func NewOrderHandler(s service.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido"})
		return
	}

	req.ID = generateOrderID()
	req.Status = models.StatusPending
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := h.service.CreateOrder(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando pedido"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	order, err := h.service.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	newStatus := models.OrderStatus(c.Query("status")) // Convertir a models.OrderStatus
	if err := h.service.UpdateOrderStatus(c.Request.Context(), orderID, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando estado del pedido"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Estado del pedido actualizado"})
}

func generateOrderID() string {
	return "ORD-" + time.Now().Format("20060102150405")
}
