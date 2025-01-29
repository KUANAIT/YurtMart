package main

import (
	"YurtMart/database"
	"YurtMart/routes"
	"YurtMart/web"
	"log"
	"net/http"
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

	routes.RegisterRoutes()

	web.SetupTemplates()

	log.Println("Server is running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
