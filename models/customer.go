package models

import (
	"YurtMart/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name"`
	Password        string             `json:"password"`
	ShippingAddress Address            `json:"shipping_address"`
	BillingAddress  Address            `json:"billing_address"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"user_cart" bson:"user_cart"`
}

func (c *Customer) HashPassword() error {
	hashedPassword, err := auth.HashPassword(c.Password)
	if err != nil {
		return err
	}
	c.Password = hashedPassword
	return nil
}

func (c *Customer) CheckPassword(providedPassword string) bool {
	return auth.CheckPassword(c.Password, providedPassword)
}
