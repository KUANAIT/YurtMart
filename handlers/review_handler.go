package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"YurtMart/sessions"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	reviewCollection, _ := database.GetCollection("YurtMart", "reviews")
	cursor, err := reviewCollection.Find(context.TODO(), bson.M{})
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
	}{
		Username:      username,
		UserID:        userID,
		Reviews:       reviews,
		CurrentUserID: currentUserID,
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

	reviewCollection, _ := database.GetCollection("YurtMart", "reviews")
	review := models.Review{
		ID:       primitive.NewObjectID(),
		UserID:   objID,
		Username: customer.Name,
		Rating:   rating,
		Text:     text,
	}

	_, err = reviewCollection.InsertOne(context.TODO(), review)
	if err != nil {
		http.Error(w, "Failed to save review", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}

func EditReview(w http.ResponseWriter, r *http.Request) {
	reviewID := r.URL.Query().Get("id")
	if reviewID == "" {
		http.Error(w, "Missing review ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		http.Error(w, "Invalid ObjectID format", http.StatusBadRequest)
		return
	}

	collection, err := database.GetCollection("YurtMart", "reviews")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var review models.Review
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/edit-review.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, review)
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	reviewID := r.FormValue("id")
	fmt.Println("Received Review ID for Delete:", reviewID)

	if reviewID == "" {
		http.Error(w, "Missing review ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		fmt.Println("Invalid ObjectID format:", reviewID) // Log error
		http.Error(w, "Invalid ObjectID format", http.StatusBadRequest)
		return
	}

	collection, err := database.GetCollection("YurtMart", "reviews")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}
