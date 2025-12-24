package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khusa-mahal/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WishlistHandler struct {
	service *services.WishlistService
}

func NewWishlistHandler(service *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		service: service,
	}
}

// GetWishlist returns the user's wishlist with enriched product details
func (h *WishlistHandler) GetWishlist(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(claims["userId"].(string))
	if err != nil {
		fmt.Printf("❌ Invalid User ID in GetWishlist: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	fmt.Printf("🔍 Fetching wishlist for user: %s\n", userID.Hex())

	_, products, err := h.service.GetWishlist(c.Context(), userID)
	if err != nil {
		fmt.Printf("❌ Failed to fetch wishlist: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch wishlist"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
	})
}

// AddItem adds an item to the wishlist
func (h *WishlistHandler) AddItem(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(claims["userId"].(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	type Request struct {
		ProductID string `json:"productId"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	fmt.Printf("❤️ Adding item %s to wishlist for user %s\n", req.ProductID, userID.Hex())

	pID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Product ID"})
	}

	if err := h.service.AddItem(c.Context(), userID, pID); err != nil {
		fmt.Printf("❌ Failed to add item: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add item"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Item added to wishlist",
	})
}

// RemoveItem removes an item from the wishlist
func (h *WishlistHandler) RemoveItem(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(claims["userId"].(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	productID := c.Params("productId")
	fmt.Printf("🗑️ Removing item %s from wishlist for user %s\n", productID, userID.Hex())

	pID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Product ID"})
	}

	if err := h.service.RemoveItem(c.Context(), userID, pID); err != nil {
		fmt.Printf("❌ Failed to remove item: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove item"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Item removed from wishlist",
	})
}
