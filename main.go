package main

import (
	"awesomeProject17/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"awesomeProject17/config"
)

func main() {
	config.ConnectDB()

	router := mux.NewRouter()
	router.HandleFunc("/items", handlers.GetItems).Methods("GET")
	router.HandleFunc("/items", handlers.CreateItem).Methods("POST")
	router.HandleFunc("/items", handlers.UpdateItem).Methods("PUT")
	router.HandleFunc("/items", handlers.DeleteItem).Methods("DELETE")

	log.Println("Server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
