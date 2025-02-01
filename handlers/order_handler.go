package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"YurtMart/database"
	"YurtMart/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddItem(w http.ResponseWriter, r *http.Request) {
	var item models.ItemOrdered
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.OrdersDB.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func Display(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.OrdersDB.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var order models.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func GetPrice(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.OrdersDB.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var totalPrice float64
	for cursor.Next(ctx) {
		var order models.Order
		cursor.Decode(&order)
		totalPrice += order.Total
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"total_price": totalPrice})
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := primitive.ObjectIDFromHex(params.Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.OrdersDB.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order deleted"})
}
