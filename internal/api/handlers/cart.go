package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// Helper to extract IDs
func (h *CartHandler) getIDs(c *fiber.Ctx) (userID string, sessionID string) {
	// 1. Try to get UserID from JWT if present (set by optional middleware or manual check)
	// We might use a middleware that sets "user" if token is valid, but doesn't block if not.
	// Or we just check the Locals set by Protected middleware if the route is protected.

	userToken := c.Locals("user")
	if userToken != nil {
		token := userToken.(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userID = claims["userId"].(string)
	}

	// 2. Get SessionID from Header
	sessionID = c.Get("X-Session-ID")
	return
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID, sessionID := h.getIDs(c)

	cart, err := h.cartService.GetCartEnriched(c.Context(), userID, sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": cart})
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID, sessionID := h.getIDs(c)
	fmt.Printf("🛒 AddToCart: UserID='%s', SessionID='%s'\n", userID, sessionID)

	// Parse request body - expecting productId as string from frontend
	var req struct {
		ProductID     string `json:"productId"`
		Quantity      int    `json:"quantity"`
		SelectedSize  string `json:"selectedSize"`
		SelectedColor string `json:"selectedColor"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Convert productId string to ObjectID
	productOID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		fmt.Printf("❌ AddToCart: Invalid ProductID '%s': %v\n", req.ProductID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Create CartItem with proper ObjectID
	item := models.CartItem{
		ProductID:     productOID,
		Quantity:      req.Quantity,
		SelectedSize:  req.SelectedSize,
		SelectedColor: req.SelectedColor,
	}

	cart, err := h.cartService.AddToCart(c.Context(), userID, sessionID, item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": cart})
}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID, sessionID := h.getIDs(c)
	fmt.Printf("🗑️  RemoveItem: UserID='%s', SessionID='%s'\n", userID, sessionID)

	var req struct {
		ProductID     string `json:"productId"`
		SelectedSize  string `json:"selectedSize"`
		SelectedColor string `json:"selectedColor"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := h.cartService.RemoveItem(c.Context(), userID, sessionID, req.ProductID, req.SelectedSize, req.SelectedColor); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Item removed"})
}

func (h *CartHandler) MergeCart(c *fiber.Ctx) error {
	// This route MUST be protected, so userID is guaranteed
	userID, sessionID := h.getIDs(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.cartService.MergeCarts(c.Context(), userID, sessionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Carts merged"})
}
