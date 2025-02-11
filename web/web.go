package web

import (
	"YurtMart/handlers"
	"html/template"
	"log"
	"net/http"
)

func SetupTemplates() {
	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/register.html")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/login.html")
	})
	http.HandleFunc("/shop", handlers.RenderItemsPage)

}
