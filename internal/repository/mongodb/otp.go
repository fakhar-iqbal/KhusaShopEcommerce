package mongodb

import (
	"context"
	"time"

	"github.com/khusa-mahal/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OTPRepository struct {
	collection *mongo.Collection
}

func NewOTPRepository(db *mongo.Database) *OTPRepository {
	collection := db.Collection("otps")

	// Create TTL index for automatic expiration
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)

	return &OTPRepository{
		collection: collection,
	}
}

func (r *OTPRepository) Save(ctx context.Context, otp *models.OTP) error {
	otp.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, otp)
	return err
}

func (r *OTPRepository) FindValidOTP(ctx context.Context, email, code string) (*models.OTP, error) {
	var otp models.OTP
	filter := bson.M{
		"email":     email,
		"code":      code,
		"expiresAt": bson.M{"$gt": time.Now()},
	}
	err := r.collection.FindOne(ctx, filter).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepository) DeleteByEmail(ctx context.Context, email string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"email": email})
	return err
}
