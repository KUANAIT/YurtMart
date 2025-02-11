package routes

import (
	"YurtMart/handlers"
	"net/http"
)

func RegisterItemOrderedRoutes() {
	http.HandleFunc("/items_ordered/add", handlers.AddItem)
	http.HandleFunc("/ordered_items/getprice/", handlers.GetPrice)
	http.HandleFunc("/items_ordered/display", handlers.Display)
}
