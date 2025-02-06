package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"PeriFyGo/config"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type VerifyController struct{}

// VerifyEmail validates the verification token and updates the user's Verified flag.
func (vc *VerifyController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusBadRequest)
		return
	}

	if claims["purpose"] != "verify" {
		http.Error(w, "Invalid token purpose", http.StatusBadRequest)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Email not found in token", http.StatusBadRequest)
		return
	}

	collection := config.DB.Database("perifygo_db").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"verified": true, "updated_at": time.Now()}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Email verified successfully!",
	})
}
