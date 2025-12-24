package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/khusa-mahal/backend/internal/config"
	"github.com/khusa-mahal/backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	config *config.CacheConfig
}

func NewCache(cfg *config.Config) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &Cache{
		client: client,
		config: &cfg.Cache,
	}
}

// Product cache operations

func (c *Cache) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	key := fmt.Sprintf("product:%s", id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *Cache) SetProduct(ctx context.Context, id string, product *models.Product) error {
	key := fmt.Sprintf("product:%s", id)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.config.ProductTTL).Err()
}

func (c *Cache) DeleteProduct(ctx context.Context, id string) error {
	key := fmt.Sprintf("product:%s", id)
	return c.client.Del(ctx, key).Err()
}

// Product list cache operations

func (c *Cache) GetProductList(ctx context.Context, filter string) ([]models.Product, error) {
	key := fmt.Sprintf("products:list:%s", filter)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err := json.Unmarshal([]byte(val), &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (c *Cache) SetProductList(ctx context.Context, filter string, products []models.Product) error {
	key := fmt.Sprintf("products:list:%s", filter)
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.config.ListTTL).Err()
}

// Cart cache operations

func (c *Cache) GetCart(ctx context.Context, sessionID string) (*models.Cart, error) {
	key := fmt.Sprintf("cart:%s", sessionID)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cart models.Cart
	if err := json.Unmarshal([]byte(val), &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *Cache) SetCart(ctx context.Context, sessionID string, cart *models.Cart) error {
	key := fmt.Sprintf("cart:%s", sessionID)
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.config.CartTTL).Err()
}

func (c *Cache) DeleteCart(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("cart:%s", sessionID)
	return c.client.Del(ctx, key).Err()
}

// User cart operations

func (c *Cache) GetUserCart(ctx context.Context, userID string) (*models.Cart, error) {
	key := fmt.Sprintf("cart:user:%s", userID)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cart models.Cart
	if err := json.Unmarshal([]byte(val), &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *Cache) SetUserCart(ctx context.Context, userID string, cart *models.Cart) error {
	key := fmt.Sprintf("cart:user:%s", userID)
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	// User carts don't expire
	return c.client.Set(ctx, key, data, 0).Err()
}

// Invalidate all product caches (useful after updates)
func (c *Cache) InvalidateProductCaches(ctx context.Context) error {
	iter := c.client.Scan(ctx, 0, "products:list:*", 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Ping checks Redis connection
func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close closes the Redis client
func (c *Cache) Close() error {
	return c.client.Close()
}
