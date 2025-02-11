package main

import (
	"YurtMart/database"
	"YurtMart/routes"
	"YurtMart/web"
	"log"
	"net/http"
)

func main() {
	database.Connect_DB()

	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := database.DisconnectDB(); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	web.SetupTemplates()
	routes.RegisterRoutes()

	routes.RegisterItemRoutes()

	log.Println("Server started on :8086")
	log.Fatal(http.ListenAndServe(":8086", nil))
}
