package web

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func SetupCORS(router *mux.Router) http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:63342"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	return corsHandler.Handler(router)
}
