package routes

import (
	"YurtMart/handlers"
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/customers", handlers.CreateCustomer)
	http.HandleFunc("/customers/get", handlers.GetCustomer)
	http.HandleFunc("/customers/update", handlers.UpdateCustomer)
	http.HandleFunc("/customers/delete", handlers.DeleteCustomer)
	http.HandleFunc("/customers/address", handlers.GetCustomerAddress)
	http.HandleFunc("/loginuser", handlers.LoginCustomer)
	http.HandleFunc("/profile", handlers.Profile)
	http.HandleFunc("/profile/edit", handlers.EditShippingAddress)
	http.HandleFunc("/profile/change-password", handlers.ChangePassword)

}
