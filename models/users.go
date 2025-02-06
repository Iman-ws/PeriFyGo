package models

import "time"

// User represents a user in the PeriFyGo system.
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"` // stored as bcrypt hash
	Role      string    `bson:"role" json:"role"`         // e.g., "admin", "user"
	Verified  bool      `bson:"verified" json:"verified"` // email verification flag
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
