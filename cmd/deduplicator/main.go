package main

import (
	"context"
	"fmt"
	"log"

	"github.com/khusa-mahal/backend/internal/config"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	log.Println("🧹 Starting product deduplication...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to MongoDB
	db, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.Close(context.Background())

	log.Println("✅ Connected to MongoDB")

	ctx := context.Background()
	collection := db.GetDB().Collection("products")

	// Find all products
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to fetch products:", err)
	}
	defer cursor.Close(ctx)

	// Map to track product names and their first occurrence
	seenProducts := make(map[string]primitive.ObjectID)
	var duplicates []primitive.ObjectID

	// Track all products
	type Product struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}

	var products []Product
	if err := cursor.All(ctx, &products); err != nil {
		log.Fatal("Failed to decode products:", err)
	}

	log.Printf("📊 Found %d total products\n", len(products))

	// Identify duplicates
	for _, product := range products {
		if firstID, exists := seenProducts[product.Name]; exists {
			// This is a duplicate - mark for deletion
			duplicates = append(duplicates, product.ID)
			log.Printf("🔍 Duplicate found: '%s' (ID: %s) - Will keep first occurrence (ID: %s)\n",
				product.Name, product.ID.Hex(), firstID.Hex())
		} else {
			// First occurrence of this product name
			seenProducts[product.Name] = product.ID
		}
	}

	if len(duplicates) == 0 {
		log.Println("✅ No duplicates found! Database is clean.")
		return
	}

	log.Printf("⚠️  Found %d duplicate products\n", len(duplicates))
	log.Println("🗑️  Deleting duplicates...")

	// Delete duplicates
	result, err := collection.DeleteMany(ctx, bson.M{
		"_id": bson.M{"$in": duplicates},
	})
	if err != nil {
		log.Fatal("Failed to delete duplicates:", err)
	}

	log.Printf("✅ Successfully deleted %d duplicate products\n", result.DeletedCount)
	log.Printf("📦 Final count: %d unique products\n", len(seenProducts))
	log.Println("🎉 Deduplication complete!")

	// Show summary
	log.Println("\n📋 Remaining unique products:")
	for name := range seenProducts {
		fmt.Printf("  ✓ %s\n", name)
	}
}
