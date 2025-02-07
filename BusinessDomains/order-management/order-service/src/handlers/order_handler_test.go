package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"order-management/order-service/src/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderService es un mock del servicio de órdenes
type MockOrderService struct {
	mock.Mock
}

// CreateOrder mocks el método CreateOrder
func (m *MockOrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// GetOrder mocks el método GetOrder
func (m *MockOrderService) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*models.Order), args.Error(1)
}

// UpdateOrderStatus mocks el método UpdateOrderStatus
func (m *MockOrderService) UpdateOrderStatus(ctx context.Context, orderID string, status models.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func TestCreateOrderHandler(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	order := &models.Order{
		ID:     "ORD-12345",
		UserID: "user1",
		Status: models.StatusPending,
		Total:  100.0,
		Items: []models.OrderItem{
			{ProductID: "prod1", Quantity: 1, Price: 50.0},
			{ProductID: "prod2", Quantity: 2, Price: 25.0},
		},
	}

	mockService.On("CreateOrder", mock.Anything, order).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/orders", nil)

	handler.CreateOrder(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetOrderHandler(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	order := &models.Order{
		ID:     "ORD-12345",
		UserID: "user1",
		Status: models.StatusPending,
		Total:  100.0,
		Items: []models.OrderItem{
			{ProductID: "prod1", Quantity: 1, Price: 50.0},
			{ProductID: "prod2", Quantity: 2, Price: 25.0},
		},
	}

	mockService.On("GetOrder", mock.Anything, "ORD-12345").Return(order, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/orders/ORD-12345", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "ORD-12345"}}

	handler.GetOrder(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateOrderStatusHandler(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	mockService.On("UpdateOrderStatus", mock.Anything, "ORD-12345", models.StatusShipped).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/orders/ORD-12345?status=SHIPPED", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "ORD-12345"}}

	handler.UpdateOrderStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
