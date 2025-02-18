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

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	// Retrieve user ID from the session
	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		// If there's no user ID in the session, render the page without user-specific data
		data := struct {
			Authenticated bool
			UserName      string
		}{
			Authenticated: false,
			UserName:      "",
		}
		tmpl.Execute(w, data)
		return
	}

	// Convert userID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Fetch user data from the database
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

	// Check if the user is authenticated
	_, isAuthenticated := session.Values["authenticated"].(bool)

	data := struct {
		Authenticated bool
		UserName      string
	}{
		Authenticated: isAuthenticated,
		UserName:      customer.Name, // Fetch the name from the database
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
