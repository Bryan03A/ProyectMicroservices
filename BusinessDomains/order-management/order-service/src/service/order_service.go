package service

import (
	"context"
	"order-management/order-service/src/models"
	"order-management/order-service/src/repository"
)

// OrderServiceInterface define los métodos que debe implementar un servicio de órdenes
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status models.OrderStatus) error
}

// OrderService es la implementación concreta del servicio de órdenes
type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Implementación del método CreateOrder
	return s.repo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	// Implementación del método GetOrder
	return s.repo.GetOrder(ctx, orderID)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, status models.OrderStatus) error {
	// Implementación del método UpdateOrderStatus
	return s.repo.UpdateOrderStatus(ctx, orderID, status)
}
