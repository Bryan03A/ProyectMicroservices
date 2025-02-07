package models

import (
	"time"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	// StatusPending indicates that the order is pending.
	StatusPending OrderStatus = "PENDING"
	// StatusConfirmed indicates that the order is confirmed.
	StatusConfirmed OrderStatus = "CONFIRMED"
	// StatusShipped indicates that the order has been shipped.
	StatusShipped OrderStatus = "SHIPPED"
	// StatusCanceled indicates that the order has been canceled.
	StatusCanceled OrderStatus = "CANCELED"
)

// Order represents an order in the system.
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Status    OrderStatus `json:"status"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// StatusUpdate represents a status update for an order.
type StatusUpdate struct {
	Status OrderStatus `json:"status"`
}
