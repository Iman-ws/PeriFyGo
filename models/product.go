package models

import "time"

// Product represents a product in the PeriFyGo store.
// It includes basic fields such as Name, Description, Price, ImageURL, and timestamps.
type Product struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB ObjectID
	Name        string    `bson:"name" json:"name"`                  // Product name
	Description string    `bson:"description" json:"description"`    // Product description
	Price       float64   `bson:"price" json:"price"`                // Product price
	ImageURL    string    `bson:"image_url" json:"image_url"`        // URL of the product image
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`      // Creation timestamp
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`      // Last update timestamp
}
