package main

import (
	"log"
	"net/http"

	"YurtMart/database"
	"YurtMart/routes"
	"YurtMart/web"
	"github.com/gorilla/mux"
)

func main() {

	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := database.DisconnectDB(); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	database.Connect_DB()

	router := mux.NewRouter()

	routes.RegisterRoutes()
	routes.RegisterItemRoutes(router)
	routes.RegisterOrderRoutes(router)

	web.SetupTemplates()

	log.Println("Server started on :8087")
	log.Fatal(http.ListenAndServe(":8087", router))
}
