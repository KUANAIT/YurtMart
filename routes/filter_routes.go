package routes

import (
	"YurtMart/handlers"
	"net/http"
)

func SearchRoutes() {
	http.HandleFunc("/items/search", handlers.GetItemsByName)
	http.HandleFunc("/items/category", handlers.GetItemsByCategory)
}
