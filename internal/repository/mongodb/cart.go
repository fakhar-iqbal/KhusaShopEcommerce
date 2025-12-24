package mongodb

import (
	"context"
	"time"

	"github.com/khusa-mahal/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository struct {
	collection *mongo.Collection
}

func NewCartRepository(db *mongo.Database) *CartRepository {
	return &CartRepository{
		collection: db.Collection("carts"),
	}
}

// FindByUserID finds a cart by the User's ID
func (r *CartRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Cart, error) {
	filter := bson.M{"userId": userID}
	var cart models.Cart
	err := r.collection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no cart found, not error
		}
		return nil, err
	}
	return &cart, nil
}

// FindBySessionID finds a cart by the Session ID
func (r *CartRepository) FindBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	filter := bson.M{"sessionId": sessionID}
	var cart models.Cart
	err := r.collection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

// Save creates or updates a cart
func (r *CartRepository) Save(ctx context.Context, cart *models.Cart) error {
	cart.UpdatedAt = time.Now()

	// If ID is missing, it's a new cart -> InsertOne
	if cart.ID == primitive.NilObjectID {
		if cart.CreatedAt.IsZero() {
			cart.CreatedAt = time.Now()
		}
		result, err := r.collection.InsertOne(ctx, cart)
		if err != nil {
			return err
		}
		// Update the struct with the new ID
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			cart.ID = oid
		}
		return nil
	}

	// If ID exists -> ReplaceOne
	filter := bson.M{"_id": cart.ID}

	// Ensure we don't accidentally overwrite CreatedAt with zero if it was somehow lost (unlikely but safe)
	if cart.CreatedAt.IsZero() {
		// Try to preserve original CreatedAt? Or just set to now if missing?
		// Ideally we fetch before save, so it should be there.
		cart.CreatedAt = time.Now()
	}

	_, err := r.collection.ReplaceOne(ctx, filter, cart)
	return err
}

// DeleteBySessionID removes a cart (e.g. after merging)
func (r *CartRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"sessionId": sessionID})
	return err
}
