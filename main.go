package main

import (
	"YurtMart/database"
	"YurtMart/handlers"
	"html/template"
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

	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // Serve static files under /static

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			http.Error(w, "Unable to load templates", http.StatusInternalServerError)
			log.Printf("Template error: %v", err)
			return
		}

		data := map[string]interface{}{
			"Title":   "YurtMart - Your Online Supermarket",
			"Message": "Discover a wide range of high-quality products at affordable prices.",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render templates", http.StatusInternalServerError)
			log.Printf("Render error: %v", err)
		}
	})

	http.HandleFunc("/customers", handlers.CreateCustomer)
	http.HandleFunc("/customers/get", handlers.GetCustomer)
	http.HandleFunc("/customers/update", handlers.UpdateCustomer)
	http.HandleFunc("/customers/delete", handlers.DeleteCustomer)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/register.html")
	})

	fs = http.FileServer(http.Dir("./templates"))

	log.Println("Server is running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
