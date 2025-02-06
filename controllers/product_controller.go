package controllers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"

	"PeriFyGo/config"
	"PeriFyGo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductController handles product-related operations.
type ProductController struct{}

// CreateProduct inserts a new product into the database.
func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	collection := config.DB.Database("perifygo_db").Collection("products")
	res, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		http.Error(w, "Error inserting product into DB", http.StatusInternalServerError)
		return
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Error parsing inserted ID", http.StatusInternalServerError)
		return
	}
	product.ID = id.Hex()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GetProducts retrieves products with basic filtering, sorting, and pagination.
func (pc *ProductController) GetProducts(w http.ResponseWriter, r *http.Request) {
	collection := config.DB.Database("perifygo_db").Collection("products")

	filterStr := r.URL.Query().Get("filter")
	sortStr := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	limit := 10
	page := 1
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	filter := bson.M{}
	if filterStr != "" {
		filter["name"] = bson.M{"$regex": filterStr, "$options": "i"}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))
	if sortStr != "" {
		order := 1
		if sortStr[0] == '-' {
			order = -1
			sortStr = sortStr[1:]
		}
		findOptions.SetSort(bson.D{{Key: sortStr, Value: order}})
	} else {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		http.Error(w, "Error reading products", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"page":     page,
		"limit":    limit,
		"products": products,
		"count":    len(products),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProduct updates an existing product by ID.
func (pc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	oid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}
	product.UpdatedAt = time.Now()

	collection := config.DB.Database("perifygo_db").Collection("products")
	update := bson.M{"$set": product}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Product updated successfully"})
}

// DeleteProduct deletes a product by ID.
func (pc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	oid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	collection := config.DB.Database("perifygo_db").Collection("products")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Product deleted successfully"})
}
