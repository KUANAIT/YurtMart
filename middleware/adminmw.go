package middleware

import (
	"YurtMart/database"
	"YurtMart/models"
	"YurtMart/sessions"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session
		session, err := sessions.Get(r)
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		// Check if the user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Authentication required"}`))
			return
		}

		// Retrieve the user ID from the session
		userID, ok := session.Values["user_id"].(string)
		if !ok || userID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "User not authenticated"}`))
			return
		}

		// Convert userID to ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Invalid user ID"}`))
			return
		}

		// Fetch the user from the database
		collection, err := database.GetCollection("YurtMart", "customers")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to get database collection"}`))
			return
		}

		var customer models.Customer
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&customer)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "User not found"}`))
			return
		}

		// Check if the user is an admin
		if !customer.Admin {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Admin access required"}`))
			return
		}

		// Call the next handler if the user is an admin
		next.ServeHTTP(w, r)
	})
}
