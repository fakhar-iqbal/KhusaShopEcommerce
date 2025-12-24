package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khusa-mahal/backend/internal/api/handlers"
	"github.com/khusa-mahal/backend/internal/api/middleware"
)

func RegisterWishlistRoutes(router fiber.Router, handler *handlers.WishlistHandler) {
	wishlist := router.Group("/wishlist")

	wishlist.Get("/", middleware.Protected(), handler.GetWishlist)
	wishlist.Post("/", middleware.Protected(), handler.AddItem)
	wishlist.Delete("/:productId", middleware.Protected(), handler.RemoveItem)
}
