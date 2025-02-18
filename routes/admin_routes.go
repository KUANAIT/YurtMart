package routes

import (
	"YurtMart/handlers" // Import the handlers package
	"YurtMart/middleware"
	"net/http"
)

func AdminRoutes() {
	http.Handle("/admin/users", middleware.AdminOnly(http.HandlerFunc(handlers.AdminListUsers)))
	http.Handle("/admin/users/delete", middleware.AdminOnly(http.HandlerFunc(handlers.AdminDeleteUser)))
	http.Handle("/admin/users/update", middleware.AdminOnly(http.HandlerFunc(handlers.AdminUpdateUser)))
	http.Handle("/admin/items", middleware.AdminOnly(http.HandlerFunc(handlers.AdminListItems)))
	http.Handle("/admin/items/delete", middleware.AdminOnly(http.HandlerFunc(handlers.AdminDeleteItem)))
	http.Handle("/admin/items/update", middleware.AdminOnly(http.HandlerFunc(handlers.AdminUpdateItem)))
	http.Handle("/admin/items/create", middleware.AdminOnly(http.HandlerFunc(handlers.AdminCreateItem)))
}
