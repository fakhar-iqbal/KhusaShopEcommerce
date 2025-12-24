package services

import (
	"context"

	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepo *mongodb.OrderRepository
}

func NewOrderService(orderRepo *mongodb.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID string) ([]models.Order, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.orderRepo.FindByUserID(ctx, oid)
}

// TODO: Implement CreateOrder when Checkout flow is ready
