package database

import (
	"YurtMart/models"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("Can't find product")
	ErrCantDecodeProducts = errors.New("Can't decode product")
	ErrUserIdIsNotValid   = errors.New("User ID is not valid")
	ErrCantUpdateUser     = errors.New("Can't update user")
	ErrCantRemoveProduct  = errors.New("Can't remove product")
	ErrCantGetItem        = errors.New("Can't get item")
	ErrCantBuyCartItem    = errors.New("Can't buy product")
)

func AddProductToCart(ctx context.Context, prodCollection, customerCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "cart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
	_, err = customerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveProductFromCart(ctx context.Context, prodCollection, customerCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"user_cart": bson.M{"_id": productID}}}
	_, err = customerCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveProduct
	}
	return nil

}

/*func BuyProductFromCart(ctx context.Context, cusomerCollection *mongo.Collection, userID string) error {
}

func InstantBuyer(){

}*/
