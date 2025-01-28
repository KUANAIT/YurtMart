package main

import (
	"YurtMart/database"
	"log"
	"net/http"
)

func main() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer func() {
		if err := database.DisconnectDB(); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Println("Server is running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
