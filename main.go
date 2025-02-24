package main

import (
	"YurtMart/database"
	"YurtMart/routes"
	"YurtMart/sessions"
	"YurtMart/web"
	"log"
	"net/http"
)

func main() {
	sessionKey := []byte("your-32-byte-secret-key-here")
	sessions.Initialize(sessionKey)

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
	routes.RegisterAuthRoutes()
	routes.RegisterItemRoutes()
	routes.RegisterItemOrderedRoutes()
	routes.SearchRoutes()
	routes.ShoppingCartRoutes()
	routes.AdminRoutes()
	routes.PaymentRoutes()
	web.SetupTemplates()

	log.Println("Server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
