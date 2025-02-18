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
	http.HandleFunc("/item", handlers.RenderItemPage)
	http.HandleFunc("/submit-review", handlers.SubmitReview)
	http.HandleFunc("/review/edit", handlers.EditReview)
	http.HandleFunc("/review/delete", handlers.DeleteReview)
	http.HandleFunc("/review", handlers.GetReview)
}
