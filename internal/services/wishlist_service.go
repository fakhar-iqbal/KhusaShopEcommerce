package services

import (
	"context"
	"fmt"

	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WishlistService struct {
	repo        *mongodb.WishlistRepository
	productRepo *mongodb.ProductRepository
}

func NewWishlistService(repo *mongodb.WishlistRepository, productRepo *mongodb.ProductRepository) *WishlistService {
	return &WishlistService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *WishlistService) GetWishlist(ctx context.Context, userID primitive.ObjectID) (*models.Wishlist, []models.Product, error) {
	wishlist, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// Enrich with products
	// Optimization: This could be better with an aggregation or a In query
	// but for simplicity looping or a Find In for now.
	// We'll trust the product repo to have a method to find multiple maybe?
	// Or we can just iterate. Given wished items are usually few (<50), it is fine.

	enrichedProducts := []models.Product{}
	// Simple approach: Get product details for valid IDs
	// For a production app we'd add FindByIDs to ProductRepo.
	// Let's assume we do loop or add FindByIDs.
	// I'll add FindByIDs to ProductRepo as a separate small edit or just FindByID in loop if lazy.
	// Loop is risky for DB hits. Let's just create a FindByIDs later or now.
	// For now, I will use a loop because I don't want to change ProductRepo interface yet if not needed.
	// Actually, mongo driver supports $in easily.

	// Changing approach: Let's assume ProductRepo doesn't have FindByIDs yet.
	// I will write this loop.
	for _, pid := range wishlist.Products {
		p, err := s.productRepo.GetByID(ctx, pid.Hex())
		if err == nil && p != nil {
			enrichedProducts = append(enrichedProducts, *p)
		} else {
			fmt.Printf("⚠️ Wishlist product not found: %s (err: %v)\n", pid.Hex(), err)
		}
	}
	fmt.Printf("✅ Wishlist enriched: %d products found for user %s\n", len(enrichedProducts), userID.Hex())

	return wishlist, enrichedProducts, nil
}

func (s *WishlistService) AddItem(ctx context.Context, userID primitive.ObjectID, productID primitive.ObjectID) error {
	return s.repo.AddItem(ctx, userID, productID)
}

func (s *WishlistService) RemoveItem(ctx context.Context, userID primitive.ObjectID, productID primitive.ObjectID) error {
	return s.repo.RemoveItem(ctx, userID, productID)
}
