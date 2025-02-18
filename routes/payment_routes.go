package routes

import (
	"YurtMart/handlers"
	"YurtMart/middleware"
	"net/http"
)

func PaymentPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/payment.html")
}

func PaymentRoutes() {
	http.HandleFunc("/payment/process", middleware.AuthRequired(handlers.ProcessPayment))
	http.HandleFunc("/payment", middleware.AuthRequired(PaymentPage))
	http.HandleFunc("/customer/shipping-address", middleware.AuthRequired(handlers.GetShippingAddress))
}
