package controllers

import (
	"YurtMart/database"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection     *mongo.Collection
	customerCollection *mongo.Collection
}

func NewApplication(prodCollection, customerCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection:     prodCollection,
		customerCollection: customerCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is missing"))
			return
		}

		userQueryID := c.Query("user_id")
		if userQueryID == "" {
			log.Println("User ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is missing"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.customerCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Product added to cart successfully")
	}

}

func (app *Application) RemoveFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is missing"))
			return
		}

		userQueryID := c.Query("user_id")
		if userQueryID == "" {
			log.Println("User ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is missing"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.RemoveProductFromCart(ctx, app.prodCollection, app.customerCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Product removed from cart successfully")

	}

}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

/*func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Panic("User ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is missing"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyProductFromCart(ctx, app.customerCollection, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON("The order placed successfully")
	}
}
func (app *Application) InstantBuyer() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is missing"))
			return
		}

		userQueryID := c.Query("user_id")
		if userQueryID == "" {
			log.Println("User ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is missing"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.customerCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Placed the order successfully")
	}

}*/
