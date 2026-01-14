package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ColorOption represents a color variant
type ColorOption struct {
	Name     string `json:"name" bson:"name"`
	Hex      string `json:"hex" bson:"hex"`
	ImageURL string `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
}

// Product represents a khusa product
type Product struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	Category         string             `json:"category" bson:"category"`
	Price            float64            `json:"price" bson:"price"`
	OriginalPrice    *float64           `json:"originalPrice,omitempty" bson:"originalPrice,omitempty"`
	Discount         *float64           `json:"discount,omitempty" bson:"discount,omitempty"`
	Image            string             `json:"image" bson:"image"`
	Images           []string           `json:"images,omitempty" bson:"images,omitempty"`
	CombinationImage string             `json:"combinationImage,omitempty" bson:"combinationImage,omitempty"`
	WornImage        string             `json:"wornImage,omitempty" bson:"wornImage,omitempty"`
	Description      string             `json:"description,omitempty" bson:"description,omitempty"`
	IsNew            bool               `json:"isNew,omitempty" bson:"isNew,omitempty"`
	IsSale           bool               `json:"isSale,omitempty" bson:"isSale,omitempty"`
	Rating           *float64           `json:"rating,omitempty" bson:"rating,omitempty"`
	Reviews          *int               `json:"reviews,omitempty" bson:"reviews,omitempty"`
	Sizes            []string           `json:"sizes" bson:"sizes"`
	Colors           []ColorOption      `json:"colors" bson:"colors"`
	Stock            int                `json:"stock" bson:"stock"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// User represents a user account
type User struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email             string             `json:"email" bson:"email"`
	PasswordHash      string             `json:"-" bson:"passwordHash"`
	Name              string             `json:"name" bson:"name"`
	Age               int                `json:"age,omitempty" bson:"age,omitempty"`
	Phone             string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Address           *Address           `json:"address,omitempty" bson:"address,omitempty"`
	IsVerified        bool               `json:"isVerified" bson:"isVerified"`
	VerificationToken string             `json:"-" bson:"verificationToken,omitempty"`
	CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// OTP represents a one-time password for email verification
type OTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	Code      string             `bson:"code"`
	ExpiresAt time.Time          `bson:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt"`
}

// Address represents a shipping/billing address
type Address struct {
	Street  string `json:"street" bson:"street"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	ZipCode string `json:"zipCode" bson:"zipCode"`
	Country string `json:"country" bson:"country"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ProductID     primitive.ObjectID `json:"productId" bson:"productId"`
	Product       *Product           `json:"product,omitempty" bson:"-"`
	Quantity      int                `json:"quantity" bson:"quantity"`
	SelectedSize  string             `json:"selectedSize" bson:"selectedSize"`
	SelectedColor string             `json:"selectedColor" bson:"selectedColor"`
}

// Cart represents a shopping cart
type Cart struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UserID    *primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	SessionID string              `json:"sessionId,omitempty" bson:"sessionId,omitempty"`
	Items     []CartItem          `json:"items" bson:"items"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt" bson:"updatedAt"`
}

// Order represents a completed order
type Order struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"userId" bson:"userId"`
	Items           []CartItem         `json:"items" bson:"items"`
	ShippingAddress Address            `json:"shippingAddress" bson:"shippingAddress"`
	SubTotal        float64            `json:"subTotal" bson:"subTotal"`
	ShippingCost    float64            `json:"shippingCost" bson:"shippingCost"`
	Total           float64            `json:"total" bson:"total"`
	PaymentMethod   string             `json:"paymentMethod" bson:"paymentMethod"` // cod, card, jazzcash, easypaisa
	Status          string             `json:"status" bson:"status"`               // pending, processing, shipped, delivered, cancelled
	PaymentStatus   string             `json:"paymentStatus" bson:"paymentStatus"` // pending, completed, failed
	SessionID       string             `json:"sessionId,omitempty" bson:"sessionId,omitempty"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// Category represents a product category
type Category struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// Review represents a product review
type Review struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductID primitive.ObjectID `json:"productId" bson:"productId"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Rating    int                `json:"rating" bson:"rating"`
	Comment   string             `json:"comment,omitempty" bson:"comment,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

// Wishlist represents a user's wishlist
type Wishlist struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID   `json:"userId" bson:"userId"`
	Products  []primitive.ObjectID `json:"products" bson:"products"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updatedAt"`
}

// OrderDetailsItem helper for emails
type OrderDetailsItem struct {
	Name     string
	Image    string
	Quantity int
	Price    float64
	Size     string
	Color    string
}
