package handlers

import (
	"YurtMart/database"
	"YurtMart/models"
	"YurtMart/sessions"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	var customer models.Customer
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&customer)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Name            string
		ID              string
		ShippingAddress models.Address
		BillingAddress  models.Address
		Admin           bool
	}{
		Name:            customer.Name,
		ID:              userID,
		ShippingAddress: customer.ShippingAddress,
		BillingAddress:  customer.BillingAddress,
		Admin:           customer.Admin,
	})
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func EditShippingAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"shippingAddress.street":     r.FormValue("street"),
			"shippingAddress.city":       r.FormValue("city"),
			"shippingAddress.state":      r.FormValue("state"),
			"shippingAddress.postalCode": r.FormValue("postal_code"),
			"shippingAddress.country":    r.FormValue("country"),
		},
	})
	if err != nil {
		http.Error(w, "Failed to update shipping address", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile?error=Invalid+request+method", http.StatusSeeOther)
		return
	}

	session, err := sessions.Get(r)
	if err != nil {
		http.Redirect(w, r, "/profile?error=Failed+to+retrieve+session", http.StatusSeeOther)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Redirect(w, r, "/profile?error=Unauthorized", http.StatusSeeOther)
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Redirect(w, r, "/profile?error=Invalid+user+ID", http.StatusSeeOther)
		return
	}

	collection, err := database.GetCollection("YurtMart", "customers")
	if err != nil {
		http.Redirect(w, r, "/profile?error=Failed+to+get+database+collection", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/profile?error=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	var customer models.Customer
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&customer)
	if err != nil {
		http.Redirect(w, r, "/profile?error=User+not+found", http.StatusSeeOther)
		return
	}

	if !customer.CheckPassword(currentPassword) {
		http.Redirect(w, r, "/profile?error=Current+password+is+incorrect", http.StatusSeeOther)
		return
	}

	if newPassword != confirmPassword {
		http.Redirect(w, r, "/profile?error=New+passwords+do+not+match", http.StatusSeeOther)
		return
	}

	customer.Password = newPassword
	err = customer.HashPassword()
	if err != nil {
		http.Redirect(w, r, "/profile?error=Failed+to+hash+new+password", http.StatusSeeOther)
		return
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{
		"$set": bson.M{"password": customer.Password},
	})
	if err != nil {
		http.Redirect(w, r, "/profile?error=Failed+to+update+password", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile?success=Password+changed+successfully", http.StatusSeeOther)
}
