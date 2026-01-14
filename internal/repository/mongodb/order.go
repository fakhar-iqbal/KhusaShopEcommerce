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

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *OrderRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.Order, error) {
	filter := bson.M{"userId": userID}
	opts := options.Find().SetSort(bson.M{"createdAt": -1}) // Newest first

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	// Ensure empty slice instead of nil if no orders
	if orders == nil {
		orders = []models.Order{}
	}

	return orders, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, orderID primitive.ObjectID) (*models.Order, error) {
	var order models.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	order.UpdatedAt = time.Now()
	filter := bson.M{"_id": order.ID}
	update := bson.M{"$set": order}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID primitive.ObjectID, status string, paymentStatus string) error {
	update := bson.M{
		"$set": bson.M{
			"status":        status,
			"paymentStatus": paymentStatus,
			"updatedAt":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": orderID}, update)
	return err
}
