package mongodb

import (
	"context"
	"time"

	"github.com/khusa-mahal/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WishlistRepository struct {
	collection *mongo.Collection
}

func NewWishlistRepository(db *mongo.Database) *WishlistRepository {
	return &WishlistRepository{
		collection: db.Collection("wishlists"),
	}
}

// GetByUserID returns the wishlist for a given user. Creates one if it doesn't exist.
func (r *WishlistRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&wishlist)
	if err == mongo.ErrNoDocuments {
		// Create new empty wishlist
		wishlist = models.Wishlist{
			UserID:    userID,
			Products:  []primitive.ObjectID{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		res, err := r.collection.InsertOne(ctx, wishlist)
		if err != nil {
			return nil, err
		}
		wishlist.ID = res.InsertedID.(primitive.ObjectID)
		return &wishlist, nil
	}
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

// AddItem adds a product ID to the user's wishlist
func (r *WishlistRepository) AddItem(ctx context.Context, userID primitive.ObjectID, productID primitive.ObjectID) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$addToSet": bson.M{"products": productID},
		"$set":      bson.M{"updatedAt": time.Now()},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// RemoveItem removes a product ID from the user's wishlist
func (r *WishlistRepository) RemoveItem(ctx context.Context, userID primitive.ObjectID, productID primitive.ObjectID) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$pull": bson.M{"products": productID},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
