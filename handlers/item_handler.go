package handlers

import (
	"YurtMart/sessions"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"

	"YurtMart/database"
	"YurtMart/models"

	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func RenderItemsPage(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item models.Item
		if err := cursor.Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	var userName string
	var isAuthenticated bool

	if ok && userID != "" {
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusInternalServerError)
			return
		}

		var customer models.Customer
		collection, err := database.GetCollection("YurtMart", "customers")
		if err != nil {
			http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
			return
		}

		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&customer)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		userName = customer.Name
		isAuthenticated = true
	}

	data := struct {
		Authenticated bool
		UserName      string
		Items         []models.Item
	}{
		Authenticated: isAuthenticated,
		UserName:      userName,
		Items:         items,
	}

	t, err := template.ParseFiles("templates/shop.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.DB.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item models.Item
		cursor.Decode(&item)
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := primitive.ObjectIDFromHex(params.Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem models.Item
	err = json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updatedItem}
	_, err = database.DB.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := primitive.ObjectIDFromHex(params.Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.DB.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
}

func RenderItemPage(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var item models.Item
	err = database.DB.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	reviewCollection, _ := database.GetCollection("YurtMart", "reviews")
	cursor, err := reviewCollection.Find(ctx, bson.M{"item_id": objID})
	if err != nil {
		http.Error(w, "Failed to retrieve reviews", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		http.Error(w, "Failed to parse reviews", http.StatusInternalServerError)
		return
	}

	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	var userName string
	var isAuthenticated bool
	var currentUserID primitive.ObjectID

	if ok && userID != "" {
		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusInternalServerError)
			return
		}

		var customer models.Customer
		collection, err := database.GetCollection("YurtMart", "customers")
		if err != nil {
			http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
			return
		}

		err = collection.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&customer)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		userName = customer.Name
		isAuthenticated = true
		currentUserID = userObjID
	}

	data := struct {
		Authenticated bool
		UserName      string
		Item          models.Item
		Reviews       []models.Review
		CurrentUserID primitive.ObjectID
	}{
		Authenticated: isAuthenticated,
		UserName:      userName,
		Item:          item,
		Reviews:       reviews,
		CurrentUserID: currentUserID,
	}

	t, err := template.ParseFiles("templates/item_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}
