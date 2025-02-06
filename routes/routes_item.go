package routes

import (
	"YurtMart/handlers"
	"github.com/gorilla/mux"
)

func SetupItemRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/items", handlers.GetItems).Methods("GET", "OPTIONS")
	router.HandleFunc("/items", handlers.CreateItem).Methods("POST", "OPTIONS")
	router.HandleFunc("/items", handlers.UpdateItem).Methods("PUT", "OPTIONS")
	router.HandleFunc("/items/delete", handlers.DeleteItem).Methods("DELETE", "OPTIONS")

	return router
}
