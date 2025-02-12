package routes

import (
	"YurtMart/handlers"
	"net/http"
)

func ShoppingCartRoutes() {
	http.HandleFunc("/cart", handlers.GetCart)
	http.HandleFunc("/cart/add", handlers.AddToCart)
	http.HandleFunc("/cart/remove", handlers.RemoveFromCart)
	http.HandleFunc("/cart/clear", handlers.ClearCart)
}
