package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"YurtMart/sessions"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, `{"error": "Session error"}`, http.StatusInternalServerError)
		return
	}

	// Check if the user is authenticated
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Convert userID to ObjectID
	customerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	// Parse the request body
	var requestData struct {
		Amount float64 `json:"amount"`
		Method string  `json:"method"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate payment details
	if requestData.Amount <= 0 || requestData.Method == "" {
		http.Error(w, `{"error": "Invalid payment details"}`, http.StatusBadRequest)
		return
	}

	// Create a new payment record
	payment := models.Payment{
		CustomerID: customerID,
		Amount:     requestData.Amount,
		Status:     "pending", // Initial status
		Method:     requestData.Method,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert the payment into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = database.PaymentCollection.InsertOne(ctx, payment)
	if err != nil {
		http.Error(w, `{"error": "Failed to process payment"}`, http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Payment processed successfully",
	})
}
