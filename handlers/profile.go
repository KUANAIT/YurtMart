package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"YurtMart/sessions"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	var customer models.Customer
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&customer)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Name            string
		ID              string
		ShippingAddress models.Address
		BillingAddress  models.Address
		Admin           bool // Add this field
	}{
		Name:            customer.Name,
		ID:              userID,
		ShippingAddress: customer.ShippingAddress,
		BillingAddress:  customer.BillingAddress,
		Admin:           customer.Admin, // Pass the Admin field
	})
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
