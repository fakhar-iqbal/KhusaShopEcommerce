package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khusa-mahal/backend/internal/api/handlers"
	"github.com/khusa-mahal/backend/internal/api/middleware"
)

func RegisterOrderRoutes(router fiber.Router, handler *handlers.OrderHandler) {
	orders := router.Group("/orders")
	// Apply JWT middleware to protect these routes
	orders.Use(middleware.Protected())

	orders.Get("/", handler.GetOrders)
	orders.Post("/", handler.CreateOrder)
}
