package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"YurtMart/sessions"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"strconv"
)

func Reviews(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	userID, _ := session.Values["user_id"].(string)

	var username string
	var currentUserID primitive.ObjectID

	if userID != "" {
		objID, err := primitive.ObjectIDFromHex(userID)
		if err == nil {
			collection, _ := database.GetCollection("YurtMart", "customers")
			var customer models.Customer
			err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&customer)
			if err == nil {
				username = customer.Name
				currentUserID = objID
			}
		}
	}

	itemID := r.URL.Query().Get("item_id") // Get the item_id from the query parameters
	var filter bson.M
	if itemID != "" {
		objID, err := primitive.ObjectIDFromHex(itemID)
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}
		filter = bson.M{"item_id": objID} // Filter reviews by item_id
	} else {
		filter = bson.M{} // No filter, fetch all reviews
	}

	reviewCollection, _ := database.GetCollection("YurtMart", "reviews")
	cursor, err := reviewCollection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to retrieve reviews", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var reviews []models.Review
	if err = cursor.All(context.TODO(), &reviews); err != nil {
		http.Error(w, "Failed to parse reviews", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./templates/reviews.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Username      string
		UserID        string
		Reviews       []models.Review
		CurrentUserID primitive.ObjectID
		ItemID        string // Add this line
	}{
		Username:      username,
		UserID:        userID,
		Reviews:       reviews,
		CurrentUserID: currentUserID,
		ItemID:        itemID, // Add this line
	})
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func SubmitReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil || rating < 1 || rating > 5 {
		http.Error(w, "Invalid rating", http.StatusBadRequest)
		return
	}
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Review text cannot be empty", http.StatusBadRequest)
		return
	}

	itemID, err := primitive.ObjectIDFromHex(r.FormValue("item_id"))
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Check if the user has already reviewed this item
	reviewCollection, _ := database.GetCollection("YurtMart", "reviews")
	filter := bson.M{"user_id": objID, "item_id": itemID}
	var existingReview models.Review
	err = reviewCollection.FindOne(context.TODO(), filter).Decode(&existingReview)
	if err == nil {
		// If a review already exists, return an error
		http.Error(w, "You have already reviewed this item", http.StatusBadRequest)
		return
	}

	// If no review exists, insert the new review
	review := models.Review{
		ID:       primitive.NewObjectID(),
		UserID:   objID,
		Username: customer.Name,
		ItemID:   itemID,
		Rating:   rating,
		Text:     text,
	}

	_, err = reviewCollection.InsertOne(context.TODO(), review)
	if err != nil {
		http.Error(w, "Failed to save review", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/item?id="+itemID.Hex(), http.StatusSeeOther)
}

func EditReview(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request method:", r.Method) // Debug log
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reviewID := r.URL.Query().Get("id")
	fmt.Println("Editing review with ID:", reviewID) // Debug log

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		fmt.Println("Invalid review ID:", err) // Debug log
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Rating int    `json:"rating"`
		Text   string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		fmt.Println("Failed to decode request body:", err) // Debug log
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Updating review with data:", requestData) // Debug log

	reviewCollection, _ := database.GetCollection("YurtMart", "reviews")
	update := bson.M{
		"$set": bson.M{
			"rating": requestData.Rating,
			"text":   requestData.Text,
		},
	}

	_, err = reviewCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		fmt.Println("Failed to update review:", err) // Debug log
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Review updated"})
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reviewID := r.URL.Query().Get("id")
	if reviewID == "" {
		http.Error(w, "Review ID is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	reviewCollection, _ := database.GetCollection("YurtMart", "reviews")
	_, err = reviewCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Review deleted"})
}

func GetReview(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the review ID from the query parameters
	reviewID := r.URL.Query().Get("id")
	fmt.Println("Fetching review with ID:", reviewID)
	if reviewID == "" {
		http.Error(w, "Review ID is required", http.StatusBadRequest)
		return
	}
	// Convert the review ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	// Fetch the review from the database
	reviewCollection, err := database.GetCollection("YurtMart", "reviews")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	var review models.Review
	err = reviewCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch review", http.StatusInternalServerError)
		}
		return
	}

	// Return the review as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}
