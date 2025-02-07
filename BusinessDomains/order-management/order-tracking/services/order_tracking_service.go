package services

import (
	"order-tracking/models"
	"errors"
)

var statuses = make(map[string]models.OrderStatus)

func TrackOrder(id string) (models.OrderStatus, error) {
	status, exists := statuses[id]
	if !exists {
		return models.OrderStatus{}, errors.New("order status not found")
	}
	return status, nil
}
