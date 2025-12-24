package mongodb

import (
	"context"
	"time"

	"github.com/khusa-mahal/backend/internal/config"
	"github.com/khusa-mahal/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

// GetAll retrieves all products with optional filtering
func (r *ProductRepository) GetAll(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]models.Product, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// GetByID retrieves a single product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetByCategory retrieves products by category
func (r *ProductRepository) GetByCategory(ctx context.Context, category string) ([]models.Product, error) {
	return r.GetAll(ctx, bson.M{"category": category})
}

// Create inserts a new product
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, id string, product *models.Product) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	product.UpdatedAt = time.Now()
	update := bson.M{"$set": product}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// Delete deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Search performs text search on products
func (r *ProductRepository) Search(ctx context.Context, query string) ([]models.Product, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"category": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	return r.GetAll(ctx, filter)
}

// CreateIndexes creates necessary indexes for optimal performance
func (r *ProductRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "price", Value: 1}}},
		{Keys: bson.D{{Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}}},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Database connection
type Database struct {
	client *mongo.Client
	db     *mongo.Database
}

func Connect(cfg *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, err
	}

	// Ping database
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(cfg.MongoDB.Database)

	return &Database{
		client: client,
		db:     db,
	}, nil
}

func (d *Database) GetDB() *mongo.Database {
	return d.db
}

func (d *Database) Close(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}
