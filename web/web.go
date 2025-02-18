package web

import (
	"YurtMart/handlers"
	"net/http"
)

func SetupTemplates() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.HomePage)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/register.html")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/login.html")
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/about.html")
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/contact.html")
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/admin_customer.html")
	})
	http.HandleFunc("/adminitem", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/admin_items.html")
	})

	http.HandleFunc("/shop", handlers.RenderItemsPage)
	http.HandleFunc("/reviews", handlers.Reviews)

}
