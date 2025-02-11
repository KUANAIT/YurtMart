package handlers

import (
	"YurtMart/sessions"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	session, _ := sessions.Get(r)
	_, isAuthenticated := session.Values["authenticated"].(bool)

	data := struct {
		Authenticated bool
	}{
		Authenticated: isAuthenticated,
	}

	tmpl.Execute(w, data)
}
