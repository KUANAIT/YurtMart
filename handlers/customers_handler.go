package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if customer.Name == "" || customer.Password == "" {
		http.Error(w, "Name and Password are required fields", http.StatusBadRequest)
		return
	}

	if err := customer.HashPassword(); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	result, err := collection.InsertOne(r.Context(), customer)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	customer.ID = result.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	var customer models.Customer
	err = collection.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&customer)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving customer", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateFields := bson.M{}

	if customer.Name != "" {
		updateFields["name"] = customer.Name
	}

	if customer.Password != "" {
		if err := customer.HashPassword(); err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		updateFields["password"] = customer.Password
	}

	if (customer.ShippingAddress != models.Address{}) {
		updateFields["shippingaddress"] = customer.ShippingAddress
	}

	if (customer.BillingAddress != models.Address{}) {
		updateFields["billingaddress"] = customer.BillingAddress
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	updateFields["updatedat"] = time.Now()

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	result, err := collection.UpdateOne(r.Context(), bson.M{"_id": objectID}, bson.M{"$set": updateFields})
	if err != nil {
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer updated successfully"})
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	result, err := collection.DeleteOne(r.Context(), bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer deleted successfully"})
}
