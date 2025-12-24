package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/khusa-mahal/backend/internal/api/handlers"
	"github.com/khusa-mahal/backend/internal/api/middleware"
	"github.com/khusa-mahal/backend/internal/api/routes"
	"github.com/khusa-mahal/backend/internal/config"
	"github.com/khusa-mahal/backend/internal/repository/elasticsearch"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"github.com/khusa-mahal/backend/internal/repository/redis"
	"github.com/khusa-mahal/backend/internal/services"
)

func main() {
	log.Println("🏁 Starting server initialization...")

	// Load configuration
	log.Println("Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Println("Configuration loaded.")

	// Connect to MongoDB
	db, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := db.Close(context.Background()); err != nil {
			log.Println("Error closing MongoDB connection:", err)
		}
	}()

	log.Println("✅ Connected to MongoDB")

	// Initialize Redis cache
	cache := redis.NewCache(cfg)
	if err := cache.Ping(context.Background()); err != nil {
		log.Println("⚠️  Redis connection failed, running without cache:", err)
	} else {
		log.Println("✅ Connected to Redis")
	}
	defer cache.Close()

	// Initialize Elasticsearch
	searchService, err := elasticsearch.NewSearchService(cfg)
	if err != nil {
		log.Println("⚠️  Elasticsearch connection failed, search may be limited:", err)
	} else {
		// Create index if it doesn't exist
		if err := searchService.CreateIndex(context.Background()); err != nil {
			log.Println("⚠️  Failed to create Elasticsearch index:", err)
		} else {
			log.Println("✅ Connected to Elasticsearch")
		}
	}

	// Initialize repositories
	productRepo := mongodb.NewProductRepository(db.GetDB())
	userRepo := mongodb.NewUserRepository(db.GetDB())
	otpRepo := mongodb.NewOTPRepository(db.GetDB())
	orderRepo := mongodb.NewOrderRepository(db.GetDB())
	cartRepo := mongodb.NewCartRepository(db.GetDB())         // [NEW]
	wishlistRepo := mongodb.NewWishlistRepository(db.GetDB()) // [NEW]

	// Initialize services
	emailService := services.NewEmailService()
	authService := services.NewAuthService(userRepo, otpRepo, emailService)
	orderService := services.NewOrderService(orderRepo)
	cartService := services.NewCartService(cartRepo, productRepo)             // [NEW]
	wishlistService := services.NewWishlistService(wishlistRepo, productRepo) // [NEW]

	// Create indexes for better performance
	if err := productRepo.CreateIndexes(context.Background()); err != nil {
		log.Println("⚠️  Failed to create MongoDB indexes:", err)
	} else {
		log.Println("✅ MongoDB indexes created")
	}

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productRepo, cache, searchService)
	authHandler := handlers.NewAuthHandler(authService)
	orderHandler := handlers.NewOrderHandler(orderService)
	cartHandler := handlers.NewCartHandler(cartService)             // [NEW]
	wishlistHandler := handlers.NewWishlistHandler(wishlistService) // [NEW]

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Khusa Mahal API",
		DisableStartupMessage: true,
	})

	// Setup middleware
	middleware.SetupMiddleware(app, cfg)

	// Setup routes
	routes.SetupRoutes(app, productHandler)
	routes.RegisterAuthRoutes(app.Group("/api/v1"), authHandler)
	routes.RegisterOrderRoutes(app.Group("/api/v1"), orderHandler)
	routes.RegisterCartRoutes(app.Group("/api/v1"), cartHandler)         // [NEW]
	routes.RegisterWishlistRoutes(app.Group("/api/v1"), wishlistHandler) // [NEW]

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("\n🛑 Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Println("Error during shutdown:", err)
		}
	}()

	// Start server
	port := cfg.Server.Port
	log.Printf("🚀 Server starting on port %s\n", port)
	log.Printf("🌐 Health check: http://localhost:%s/api/v1/health\n", port)
	log.Printf("📦 Products API: http://localhost:%s/api/v1/products\n", port)
	log.Printf("🔐 Auth API: http://localhost:%s/api/v1/auth\n", port)
	log.Printf("📦 Orders API: http://localhost:%s/api/v1/orders\n", port)
	log.Printf("🛒 Cart API: http://localhost:%s/api/v1/cart\n", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
