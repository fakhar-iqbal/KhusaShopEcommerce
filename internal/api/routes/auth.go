package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khusa-mahal/backend/internal/api/handlers"
)

func RegisterAuthRoutes(router fiber.Router, handler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/verify-otp", handler.VerifyOTP)
	auth.Post("/login", handler.Login)
}
