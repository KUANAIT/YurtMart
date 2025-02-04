package main

import (
	"log"
	"net/http"

	"YurtMart/database"
	"YurtMart/routes"
	"YurtMart/web"
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

	routes.RegisterRoutes()
	routes.RegisterItemRoutes()
	routes.RegisterOrderRoutes()

	web.SetupTemplates()

	log.Println("Server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))

}
