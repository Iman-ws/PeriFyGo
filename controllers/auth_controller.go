package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"PeriFyGo/config"
	"PeriFyGo/models"
	"PeriFyGo/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthController handles user registration and login.
type AuthController struct{}

// Helper function to respond with a JSON error.
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Register registers a new user, hashes the password, sets Verified=false,
// generates an email verification token, and sends a verification email.
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if user.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password cannot be empty")
		return
	}

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = hashedPass

	if user.Role == "" {
		user.Role = "user"
	}
	user.Verified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := config.DB.Database("perifygo_db").Collection("users")
	count, err := collection.CountDocuments(context.Background(), bson.M{"email": user.Email})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking existing email")
		return
	}
	if count > 0 {
		respondWithError(w, http.StatusConflict, "Email already registered")
		return
	}

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error inserting user into DB")
		return
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Error parsing inserted ID")
		return
	}
	log.Printf("[DEBUG] Register: user %s created with _id=%s\n", user.Email, insertedID.Hex())

	// Generate email verification token (expires in 1 hour)
	verifyToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"purpose": "verify",
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"iat":     time.Now().Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := verifyToken.SignedString([]byte(secret))
	if err != nil {
		log.Printf("[ERROR] Failed to generate verification token: %v", err)
	}

	verifyLink := fmt.Sprintf("http://%s/verify?token=%s", r.Host, tokenString)
	emailSubject := "PeriFyGo Email Verification"
	emailBody := fmt.Sprintf("Hello %s,\n\nPlease verify your email by clicking the link below:\n%s\n\nThank you!", user.Name, verifyLink)

	if err := utils.SendEmail(user.Email, emailSubject, emailBody); err != nil {
		log.Printf("[ERROR] Failed to send verification email to %s: %v", user.Email, err)
	} else {
		log.Printf("[DEBUG] Verification email sent to %s", user.Email)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"user_id": insertedID.Hex(),
		"message": "Registration successful! Please check your email to verify your account.",
	})
}

// Login authenticates the user, verifies password and email verification,
// then issues a JWT token.
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("[DEBUG] Login attempt: email=%s, password=%s", creds.Email, creds.Password)

	collection := config.DB.Database("perifygo_db").Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": creds.Email}).Decode(&user)
	if err != nil {
		_ = utils.CheckPasswordHash(creds.Password, "$2a$12$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		log.Printf("[DEBUG] Login failed: user not found for email=%s", creds.Email)
		respondWithError(w, http.StatusUnauthorized, "Wrong credentials")
		return
	}

	log.Printf("[DEBUG] Found user in DB: email=%s, hashed=%s", user.Email, user.Password)
	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		log.Printf("[DEBUG] Password verification failed for email=%s. Provided=%s, StoredHash=%s", creds.Email, creds.Password, user.Password)
		respondWithError(w, http.StatusUnauthorized, "Wrong credentials")
		return
	}

	if !user.Verified {
		respondWithError(w, http.StatusForbidden, "Email not verified. Please check your email.")
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Role)
	if err != nil {
		log.Printf("[ERROR] Failed to generate token for user=%s: %v", user.Email, err)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Printf("[DEBUG] Login success for email=%s, generated JWT", user.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token":  token,
	})
}
