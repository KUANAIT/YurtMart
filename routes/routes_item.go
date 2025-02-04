package routes

import (
	"YurtMart/handlers"
	"github.com/gorilla/mux"
)

func RegisterItemRoutes(router *mux.Router) {

	router.HandleFunc("/items", handlers.GetItems).Methods("GET")

	router.HandleFunc("/items", handlers.CreateItem).Methods("POST")
	router.HandleFunc("/items", handlers.UpdateItem).Methods("PUT")
	router.HandleFunc("/items/delete", handlers.DeleteItem).Methods("DELETE")
}
