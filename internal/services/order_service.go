package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepo      *mongodb.OrderRepository
	paymentService *PaymentService
	emailService   *EmailService
	userRepo       *mongodb.UserRepository    // To get user email
	productRepo    *mongodb.ProductRepository // To get product details for email
}

func NewOrderService(orderRepo *mongodb.OrderRepository, paymentService *PaymentService, emailService *EmailService, userRepo *mongodb.UserRepository, productRepo *mongodb.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		paymentService: paymentService,
		emailService:   emailService,
		userRepo:       userRepo,
		productRepo:    productRepo,
	}
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID string) ([]models.Order, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.orderRepo.FindByUserID(ctx, oid)
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []models.CartItem, shippingAddress models.Address, paymentMethod string, paymentDetails map[string]interface{}, subTotal, shippingCost, total float64) (*models.Order, error) {
	// 1. Validate inputs (simplified)
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// 2. Process Payment
	paymentResult, err := s.paymentService.ProcessPayment(total, "PKR", paymentMethod, paymentDetails)
	if err != nil {
		return nil, err
	}

	if !paymentResult.Success {
		return nil, errors.New("payment failed: " + paymentResult.Message)
	}

	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	order := &models.Order{
		ID:              primitive.NewObjectID(),
		UserID:          userOID,
		Items:           items,
		ShippingAddress: shippingAddress,
		SubTotal:        subTotal,
		ShippingCost:    shippingCost,
		Total:           total,
		PaymentMethod:   paymentMethod,
		Status:          "pending", // Initial status
		PaymentStatus:   paymentResult.Status,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// 3. Send Email (Async)
	go func() {
		fmt.Printf(" [DEBUG] Starting async email process for OrderID: %s, UserID: %s\n", order.ID.Hex(), userID)
		// Fetch user to get email
		user, err := s.userRepo.FindByID(context.Background(), userOID)
		if err != nil {
			fmt.Printf(" [ERROR] Failed to find user for email: %v\n", err)
			return
		}
		if user == nil {
			fmt.Printf(" [ERROR] User not found (nil) for ID: %s\n", userID)
			return
		}

		fmt.Printf(" [DEBUG] User found: %s. Sending to email: %s\n", user.Name, user.Email)

		var emailItems []models.OrderDetailsItem
		for _, item := range items {
			pName := "Product"
			var price float64
			var imageURL string

			if item.Product != nil {
				pName = item.Product.Name
				price = item.Product.Price
				imageURL = item.Product.Image
			} else {
				// Product not populated, fetch from database
				prod, err := s.productRepo.GetByID(context.Background(), item.ProductID.Hex())
				if err == nil && prod != nil {
					pName = prod.Name
					price = prod.Price
					imageURL = prod.Image
				} else {
					fmt.Printf(" [WARN] Failed to fetch product for ID %s: %v\n", item.ProductID.Hex(), err)
					price = 0
					imageURL = "https://via.placeholder.com/80"
				}
			}

			emailItems = append(emailItems, models.OrderDetailsItem{
				Name:     pName,
				Image:    imageURL,
				Quantity: item.Quantity,
				Price:    price,
				Size:     item.SelectedSize,
				Color:    item.SelectedColor,
			})
		}

		err = s.emailService.SendOrderConfirmationEmail(user.Email, order.ID.Hex(), emailItems, shippingAddress, total)
		if err != nil {
			fmt.Printf(" [ERROR] EmailService returned error: %v\n", err)
		} else {
			fmt.Printf(" [SUCCESS] Email sent successfully to %s\n", user.Email)
		}
	}()

	return order, nil
}
