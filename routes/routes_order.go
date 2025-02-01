package routes

import (
	"YurtMart/handlers"
	"github.com/gorilla/mux"
)

func RegisterOrderRoutes(router *mux.Router) {
	router.HandleFunc("/orders/add", handlers.AddItem).Methods("POST")
	router.HandleFunc("/orders/display", handlers.Display).Methods("GET")
	router.HandleFunc("/orders/getprice", handlers.GetPrice).Methods("GET")
	router.HandleFunc("/orders/delete", handlers.DeleteOrder).Methods("DELETE")
}
