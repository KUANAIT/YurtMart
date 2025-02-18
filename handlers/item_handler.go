package handlers

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
	"strings"

	"YurtMart/database"
	"YurtMart/models"

	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func RenderItemsPage(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"category", 1}, {"name", 1}})

	cursor, err := database.DB.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	categoryMap := make(map[string][]models.Item)
	for cursor.Next(ctx) {
		var item models.Item
		if err := cursor.Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categoryMap[item.Category] = append(categoryMap[item.Category], item)
	}

	t, err := template.ParseFiles("templates/shop.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, categoryMap)
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
func UpdateItemByID(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/items/update/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updatedItem}
	_, err = database.DB.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/items/delete/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.DB.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
}
