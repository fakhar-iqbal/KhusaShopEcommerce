package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService struct {
	cartRepo    *mongodb.CartRepository
	productRepo *mongodb.ProductRepository
}

func NewCartService(cartRepo *mongodb.CartRepository, productRepo *mongodb.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetCartEnriched retrieves cart and enriches items with product details
func (s *CartService) GetCartEnriched(ctx context.Context, userID string, sessionID string) (*models.Cart, error) {
	fmt.Printf("🛒 GetCartEnriched: Called with UserID='%s', SessionID='%s'\n", userID, sessionID)

	cart, err := s.GetCart(ctx, userID, sessionID)
	if err != nil {
		fmt.Printf("❌ GetCartEnriched: GetCart failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("🛒 GetCartEnriched: Retrieved cart with %d items\n", len(cart.Items))

	// Enrich cart items with product details
	for i := range cart.Items {
		fmt.Printf("🔍 GetCartEnriched: Enriching item %d, ProductID=%s\n", i, cart.Items[i].ProductID.Hex())
		product, err := s.productRepo.GetByID(ctx, cart.Items[i].ProductID.Hex())
		if err != nil {
			fmt.Printf("⚠️ GetCartEnriched: Product not found for ID %s: %v\n", cart.Items[i].ProductID.Hex(), err)
			continue
		}
		if product != nil {
			cart.Items[i].Product = product
			fmt.Printf("✅ GetCartEnriched: Enriched item %d with product '%s'\n", i, product.Name)
		}
	}

	fmt.Printf("✅ GetCartEnriched: Returning cart with %d enriched items\n", len(cart.Items))
	return cart, nil
}

// GetCart retrieves cart based on UserID (if present) or SessionID
func (s *CartService) GetCart(ctx context.Context, userID string, sessionID string) (*models.Cart, error) {
	if userID != "" {
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		cart, err := s.cartRepo.FindByUserID(ctx, oid)
		if err != nil {
			return nil, err
		}
		if cart == nil {
			// Return empty cart structure if not found (lazy creation on add, or return empty now)
			return &models.Cart{UserID: &oid, Items: []models.CartItem{}}, nil
		}
		return cart, nil
	}

	if sessionID != "" {
		cart, err := s.cartRepo.FindBySessionID(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		if cart == nil {
			return &models.Cart{SessionID: sessionID, Items: []models.CartItem{}}, nil
		}
		return cart, nil
	}

	return nil, errors.New("no user id or session id provided")
}

// AddToCart adds an item or updates quantity if exists
func (s *CartService) AddToCart(ctx context.Context, userID, sessionID string, item models.CartItem) (*models.Cart, error) {
	fmt.Printf("🛒 Service.AddToCart: Starting for UserID='%s', SessionID='%s', ProductID='%s'\n", userID, sessionID, item.ProductID)

	// 1. Get existing cart or create new model
	cart, err := s.GetCart(ctx, userID, sessionID)
	if err != nil {
		fmt.Printf("❌ Service.AddToCart: GetCart failed: %v\n", err)
		return nil, err
	}
	fmt.Printf("🛒 Service.AddToCart: Cart found/created. Existing Items: %d\n", len(cart.Items))

	// 2. Validate Product (optional but good practice)
	// product, err := s.productRepo.GetByID(...)

	// 3. Update Items
	found := false
	for i, existing := range cart.Items {
		if existing.ProductID == item.ProductID && existing.SelectedSize == item.SelectedSize && existing.SelectedColor == item.SelectedColor {
			cart.Items[i].Quantity += item.Quantity
			found = true
			fmt.Printf("🛒 Service.AddToCart: Updated quantity for existing item\n")
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, item)
		fmt.Printf("🛒 Service.AddToCart: Appended new item\n")
	}

	// 4. Save
	if userID != "" {
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			fmt.Printf("❌ Service.AddToCart: Invalid UserID Hex: %v\n", err)
			return nil, err
		}
		cart.UserID = &oid
		cart.SessionID = "" // Clear session ID if user is logged in
	} else {
		cart.SessionID = sessionID
	}

	if err := s.cartRepo.Save(ctx, cart); err != nil {
		fmt.Printf("❌ Service.AddToCart: Save failed: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ Service.AddToCart: Save successful\n")

	return cart, nil
}

// RemoveItem removes a specific item from the cart
func (s *CartService) RemoveItem(ctx context.Context, userID, sessionID, productID, size, color string) error {
	cart, err := s.GetCart(ctx, userID, sessionID)
	if err != nil || cart == nil {
		return err
	}

	// Convert productID to ObjectID for comparison
	oid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err // Invalid ID format
	}

	newItems := []models.CartItem{}
	for _, item := range cart.Items {
		if !(item.ProductID == oid && item.SelectedSize == size && item.SelectedColor == color) {
			newItems = append(newItems, item)
		}
	}
	cart.Items = newItems

	// Save changes
	if userID != "" {
		oid, _ := primitive.ObjectIDFromHex(userID)
		cart.UserID = &oid
	} else {
		cart.SessionID = sessionID
	}
	return s.cartRepo.Save(ctx, cart)
}

// MergeCarts moves items from session cart to user cart
func (s *CartService) MergeCarts(ctx context.Context, userID, sessionID string) error {
	if userID == "" || sessionID == "" {
		return nil // Nothing to do
	}

	sessionCart, err := s.cartRepo.FindBySessionID(ctx, sessionID)
	if err != nil || sessionCart == nil || len(sessionCart.Items) == 0 {
		return nil // Nothing to merge
	}

	// Get User Cart
	userOID, _ := primitive.ObjectIDFromHex(userID)
	userCart, err := s.cartRepo.FindByUserID(ctx, userOID)
	if err != nil {
		return err
	}
	if userCart == nil {
		userCart = &models.Cart{UserID: &userOID, Items: []models.CartItem{}}
	}

	// Merge Logic
	for _, sessionItem := range sessionCart.Items {
		found := false
		for i, userItem := range userCart.Items {
			if userItem.ProductID == sessionItem.ProductID && userItem.SelectedSize == sessionItem.SelectedSize && userItem.SelectedColor == sessionItem.SelectedColor {
				userCart.Items[i].Quantity += sessionItem.Quantity
				found = true
				break
			}
		}
		if !found {
			userCart.Items = append(userCart.Items, sessionItem)
		}
	}

	// Save User Cart
	if err := s.cartRepo.Save(ctx, userCart); err != nil {
		return err
	}

	// Delete Session Cart
	return s.cartRepo.DeleteBySessionID(ctx, sessionID)
}

// TODO: Helper to Attach Product Details (Images/Names) for frontend display if needed
// For now, frontend often has product data or we can expand later
