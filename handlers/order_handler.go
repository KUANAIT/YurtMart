package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

func AddItem(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ItemsOrdered []models.ItemOrdered `json:"items"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var totalPrice float64
	var orderItems []models.ItemOrdered

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, itemOrdered := range request.ItemsOrdered {
		var item models.Item
		err = database.DB.FindOne(ctx, bson.M{"_id": itemOrdered.ItemID}).Decode(&item)
		if err != nil {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		itemOrdered.Price = item.Price

		orderItems = append(orderItems, itemOrdered)
		totalPrice += item.Price * float64(itemOrdered.Quantity)
	}

	order := models.Order{
		ID:         primitive.NewObjectID(),
		TotalPrice: totalPrice,
		Items:      orderItems,
	}

	_, err = database.ItemsOrderedCollection.InsertOne(ctx, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func GetPrice(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/ordered_items/getprice/"):]

	if idParam == "" {
		http.Error(w, "Missing order_id", http.StatusBadRequest)
		return
	}

	orderID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid order_id format", http.StatusBadRequest)
		return
	}

	var order models.Order
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Looking for order with _id: %v", orderID)

	err = database.ItemsOrderedCollection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"total_price": order.TotalPrice})
}

func Display(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.ItemsOrderedCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var orders []struct {
		ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

		TotalPrice float64 `bson:"total_price" json:"total_price"`
		Items      []struct {
			ItemID   primitive.ObjectID `bson:"item_id" json:"item_id"`
			Quantity int                `bson:"quantity" json:"quantity"`
			Price    float64            `bson:"price" json:"price"`
			Name     string             `bson:"name" json:"name"`
		} `bson:"items" json:"items"`
	}

	for cursor.Next(ctx) {
		var order struct {
			ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

			TotalPrice float64 `bson:"total_price" json:"total_price"`
			Items      []struct {
				ItemID   primitive.ObjectID `bson:"item_id" json:"item_id"`
				Quantity int                `bson:"quantity" json:"quantity"`
				Price    float64            `bson:"price" json:"price"`
				Name     string             `bson:"name" json:"name"`
			} `bson:"items" json:"items"`
		}
		if err := cursor.Decode(&order); err != nil {
			http.Error(w, "Failed to decode order", http.StatusInternalServerError)
			return
		}

		for i, item := range order.Items {
			var product models.Item
			err := database.DB.FindOne(ctx, bson.M{"_id": item.ItemID}).Decode(&product)
			if err != nil {
				continue
			}
			order.Items[i].Name = product.Name
		}

		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Cursor error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
