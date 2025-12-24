package routes

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khusa-mahal/backend/internal/api/handlers"
	"github.com/khusa-mahal/backend/internal/api/middleware"
)

// OptionalAuth middleware attempts to parse token but doesn't error if missing
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				secret := os.Getenv("JWT_SECRET")
				if secret == "" {
					secret = "change-this-secret"
				}
				return []byte(secret), nil
			})

			if err != nil {
				// Log error for debugging if needed, but OptionalAuth should just fail to Guest
				// fmt.Println("OptionalAuth Error:", err)
			}

			if token != nil && token.Valid {
				c.Locals("user", token)
			}
		}
		return c.Next()
	}
}

func RegisterCartRoutes(router fiber.Router, handler *handlers.CartHandler) {
	cart := router.Group("/cart")

	// Get and Add can be Guest or User
	cart.Get("/", OptionalAuth(), handler.GetCart)
	cart.Post("/", OptionalAuth(), handler.AddToCart)
	cart.Delete("/item", OptionalAuth(), handler.RemoveItem)

	// Merge MUST be User
	cart.Post("/merge", middleware.Protected(), handler.MergeCart)
}
