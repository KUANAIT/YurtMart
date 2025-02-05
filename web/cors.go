package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

// CORSHandler возвращает обработчик для настройки CORS
func CORSHandler(router *mux.Router) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:63342"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)(router)
}
