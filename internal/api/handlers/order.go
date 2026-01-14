package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type OrderRequestItem struct {
	ID            string `json:"id"`
	ProductID     string `json:"productId"`
	Quantity      int    `json:"quantity"`
	SelectedSize  string `json:"selectedSize"`
	SelectedColor string `json:"selectedColor"`
}

type PaymentDetailsRequest struct {
	WalletPhone string `json:"walletPhone"`
}

type CreateOrderRequest struct {
	Items           []OrderRequestItem    `json:"items"`
	ShippingAddress models.Address        `json:"shippingAddress"`
	PaymentMethod   string                `json:"paymentMethod"`
	PaymentDetails  PaymentDetailsRequest `json:"paymentDetails"`
	SubTotal        float64               `json:"subTotal"`
	ShippingCost    float64               `json:"shippingCost"`
	Total           float64               `json:"total"`
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["userId"].(string)

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		// Log the error for debugging
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
	}

	// Transform/Validate Items
	var cartItems []models.CartItem
	for _, item := range req.Items {
		idStr := item.ProductID
		if idStr == "" {
			idStr = item.ID // Fallback to "id"
		}

		pid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue // Skip invalid items or return error? Skipping for now to avoid crash
		}

		cartItems = append(cartItems, models.CartItem{
			ProductID:     pid,
			Quantity:      item.Quantity,
			SelectedSize:  item.SelectedSize,
			SelectedColor: item.SelectedColor,
		})
	}

	if len(cartItems) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No valid items in order"})
	}

	paymentDetails := map[string]interface{}{
		"walletPhone": req.PaymentDetails.WalletPhone,
	}

	order, err := h.orderService.CreateOrder(
		c.Context(),
		userID,
		cartItems,
		req.ShippingAddress,
		req.PaymentMethod,
		paymentDetails,
		req.SubTotal,
		req.ShippingCost,
		req.Total,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    order,
	})
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	// 1. Get User ID from JWT (set by middleware in Locals)
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["userId"].(string)

	// 2. Fetch Orders
	orders, err := h.orderService.GetUserOrders(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    orders,
	})
}
