package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/khusa-mahal/backend/internal/repository/elasticsearch"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"github.com/khusa-mahal/backend/internal/repository/redis"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductHandler struct {
	repo   *mongodb.ProductRepository
	cache  *redis.Cache
	search *elasticsearch.SearchService
}

func NewProductHandler(repo *mongodb.ProductRepository, cache *redis.Cache, search *elasticsearch.SearchService) *ProductHandler {
	return &ProductHandler{
		repo:   repo,
		cache:  cache,
		search: search,
	}
}

// GetProducts retrieves all products with optional filters
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	ctx := context.Background()
	category := c.Query("category")
	filter := "all"

	if category != "" {
		filter = category
	}

	// Try cache first
	products, err := h.cache.GetProductList(ctx, filter)
	if err == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    products,
			"cached":  true,
		})
	}

	// Cache miss - fetch from database
	var dbFilter bson.M
	if category != "" {
		dbFilter = bson.M{"category": category}
	} else {
		dbFilter = bson.M{}
	}

	products, err = h.repo.GetAll(ctx, dbFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch products",
		})
	}

	// Cache the result
	_ = h.cache.SetProductList(ctx, filter, products)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
		"cached":  false,
	})
}

// GetProduct retrieves a single product by ID
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	// Try cache first
	product, err := h.cache.GetProduct(ctx, id)
	if err == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    product,
			"cached":  true,
		})
	}

	// Cache miss - fetch from database
	product, err = h.repo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Product not found",
		})
	}

	// Cache the result
	_ = h.cache.SetProduct(ctx, id, product)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    product,
		"cached":  false,
	})
}

// SearchProducts searches products using Elasticsearch
func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	ctx := context.Background()
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Search query is required",
		})
	}

	products, err := h.search.SearchProducts(ctx, query, 0, 50)
	if err != nil {
		// Fallback to MongoDB search
		products, err = h.repo.Search(ctx, query)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   "Search failed",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
	})
}
