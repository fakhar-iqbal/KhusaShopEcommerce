package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khusa-mahal/backend/internal/api/handlers"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, productHandler *handlers.ProductHandler) {
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Product routes
	products := api.Group("/products")
	products.Get("/", productHandler.GetProducts)
	products.Get("/search", productHandler.SearchProducts)
	products.Get("/:id", productHandler.GetProduct)

	// TODO: Add more routes
	// - Cart routes
	// - Order routes
	// - Auth routes
	// - Category routes
}
