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
	"log"
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
		log.Printf("Cart not found for user %s, creating a new one", customerID.Hex())
		cart = models.ShoppingCart{
			ID:         primitive.NewObjectID(),
			CustomerID: customerID,
			Items:      []models.ItemOrdered{},
			TotalPrice: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	for i, item := range cart.Items {
		var product models.Item
		err := database.DB.FindOne(ctx, bson.M{"_id": item.ItemID}).Decode(&product)
		if err != nil {
			log.Printf("Error fetching product details for item %s: %v", item.ItemID.Hex(), err)
			continue
		}
		cart.Items[i].Name = product.Name
	}

	log.Printf("Returning cart data: %+v", cart)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, `{"error": "Session error"}`, http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	customerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	var requestData struct {
		ItemID   string  `json:"item_id"` // Frontend sends item_id as a string
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	itemID, err := primitive.ObjectIDFromHex(requestData.ItemID)
	if err != nil {
		http.Error(w, `{"error": "Invalid item ID"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": customerID}
	var cart models.ShoppingCart
	err = database.ShoppingCartCollection.FindOne(ctx, filter).Decode(&cart)

	if err != nil {
		cart = models.ShoppingCart{
			ID:         primitive.NewObjectID(),
			CustomerID: customerID,
			Items: []models.ItemOrdered{
				{
					ItemID:   itemID,
					Quantity: requestData.Quantity,
					Price:    requestData.Price,
				},
			},
			TotalPrice: requestData.Price * float64(requestData.Quantity),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err = database.ShoppingCartCollection.InsertOne(ctx, cart)
	} else {
		itemIndex := -1
		for i, cartItem := range cart.Items {
			if cartItem.ItemID == itemID {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			cart.Items[itemIndex].Quantity += requestData.Quantity
		} else {
			cart.Items = append(cart.Items, models.ItemOrdered{
				ItemID:   itemID,
				Quantity: requestData.Quantity,
				Price:    requestData.Price,
			})
		}

		cart.TotalPrice = 0
		for _, item := range cart.Items {
			cart.TotalPrice += item.Price * float64(item.Quantity)
		}

		update := bson.M{
			"$set": bson.M{
				"items":       cart.Items,
				"total_price": cart.TotalPrice,
				"updated_at":  time.Now(),
			},
		}
		_, err = database.ShoppingCartCollection.UpdateOne(ctx, filter, update)
	}

	if err != nil {
		log.Printf("Error updating cart: %v", err)
		http.Error(w, `{"error": "Could not add item to cart"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart"})
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

func UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, `{"error": "Session error"}`, http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	customerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	var requestData struct {
		ItemID   string `json:"item_id"` // Frontend sends item_id as a string
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if requestData.Quantity < 0 {
		http.Error(w, `{"error": "Quantity cannot be negative"}`, http.StatusBadRequest)
		return
	}

	itemID, err := primitive.ObjectIDFromHex(requestData.ItemID)
	if err != nil {
		http.Error(w, `{"error": "Invalid item ID"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": customerID}
	var cart models.ShoppingCart
	err = database.ShoppingCartCollection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		http.Error(w, `{"error": "Cart not found"}`, http.StatusNotFound)
		return
	}

	itemIndex := -1
	for i, cartItem := range cart.Items {
		if cartItem.ItemID == itemID {
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		http.Error(w, `{"error": "Item not found in cart"}`, http.StatusNotFound)
		return
	}

	cart.Items[itemIndex].Quantity = requestData.Quantity

	cart.TotalPrice = 0
	for _, item := range cart.Items {
		cart.TotalPrice += item.Price * float64(item.Quantity)
	}

	update := bson.M{
		"$set": bson.M{
			"items":       cart.Items,
			"total_price": cart.TotalPrice,
			"updated_at":  time.Now(),
		},
	}
	_, err = database.ShoppingCartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating cart: %v", err)
		http.Error(w, `{"error": "Could not update cart"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart updated"})
}
