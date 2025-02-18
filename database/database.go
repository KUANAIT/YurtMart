package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Collection
var ItemsOrderedCollection *mongo.Collection
var ShoppingCartCollection *mongo.Collection
var CustomerCollection *mongo.Collection
var PaymentCollection *mongo.Collection

func Connect_DB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	DB = client.Database("supermarket").Collection("items")
	CustomerCollection = client.Database("YurtMart").Collection("customers")
	ItemsOrderedCollection = client.Database("supermarket").Collection("ordered_items")
	ShoppingCartCollection = client.Database("supermarket").Collection("shopping_carts")
	PaymentCollection = client.Database("YurtMart").Collection("payments")
}
