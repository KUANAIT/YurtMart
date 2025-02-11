package routes

import (
	"YurtMart/handlers"
	"YurtMart/middleware"
	"net/http"
)

func RegisterAuthRoutes() {
	logoutHandler := middleware.AuthRequired(http.HandlerFunc(handlers.LogoutCustomer))
	http.Handle("/logout", logoutHandler)
}
