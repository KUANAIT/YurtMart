package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func GetCustomerAddress(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
	address := map[string]interface{}{
		"shipping_address": customer.ShippingAddress,
	}

	json.NewEncoder(w).Encode(address)
}
