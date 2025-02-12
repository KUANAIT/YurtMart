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
	"net/http"
	"time"
)

func GetCart(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	customerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.ShoppingCart
	filter := bson.M{"customer_id": customerID}
	err = database.ShoppingCartCollection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		cart = models.ShoppingCart{
			ID:         primitive.NewObjectID(),
			CustomerID: customerID,
			Items:      []models.ItemOrdered{},
			TotalPrice: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	customerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var item models.ItemOrdered
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.ShoppingCart
	filter := bson.M{"customer_id": customerID}
	err = database.ShoppingCartCollection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		cart = models.ShoppingCart{
			ID:         primitive.NewObjectID(),
			CustomerID: customerID,
			Items:      []models.ItemOrdered{item},
			TotalPrice: item.Price * float64(item.Quantity),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err = database.ShoppingCartCollection.InsertOne(ctx, cart)
	} else {
		cart.Items = append(cart.Items, item)
		cart.TotalPrice += item.Price * float64(item.Quantity)
		cart.UpdatedAt = time.Now()
		_, err = database.ShoppingCartCollection.UpdateOne(ctx, filter, bson.M{"$set": cart})
	}

	if err != nil {
		http.Error(w, "Could not add item to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Item added to cart")
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	customerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var item models.ItemOrdered
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": customerID}
	update := bson.M{"$pull": bson.M{"items": bson.M{"item_id": item.ItemID}}}

	_, err = database.ShoppingCartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Could not remove item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Item removed from cart")
}

// ClearCart empties the shopping cart
func ClearCart(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	customerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": customerID}
	update := bson.M{"$set": bson.M{"items": []models.ItemOrdered{}, "total_price": 0}}

	_, err = database.ShoppingCartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Could not clear cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Cart cleared")
}
