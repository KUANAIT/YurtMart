package routes

import (
	"YurtMart/handlers"
	"net/http"
)

func RegisterItemRoutes() {

	http.HandleFunc("/items/get", handlers.GetItems)
	http.HandleFunc("/items/create", handlers.CreateItem)
	http.HandleFunc("/items/update", handlers.UpdateItem)
	http.HandleFunc("/items/delete", handlers.DeleteItem)
}
