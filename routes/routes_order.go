package routes

import (
	"YurtMart/handlers"
	"net/http"
)

func RegisterOrderRoutes() {
	http.HandleFunc("/orders/add", handlers.AddItem)
	http.HandleFunc("/orders/display", handlers.Display)
	http.HandleFunc("/orders/getprice", handlers.GetPrice)
	http.HandleFunc("/orders/delete", handlers.DeleteOrder)
}
