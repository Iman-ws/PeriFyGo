package controllers

import (
	"PeriFyGo/config"
	"PeriFyGo/models"
	"PeriFyGo/utils"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// ProfileController handles user profile-related requests.
type ProfileController struct{}

// GetProfile returns the authenticated user's profile data.
// It expects the JWT token in the Authorization header.
func (pc *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}
	// Assume "Bearer <token>"
	tokenStr := authHeader[len("Bearer "):]
	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	collection := config.DB.Database("perifygo_db").Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"email": claims.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Remove password from output.
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
