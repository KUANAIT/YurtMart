package models

import (
	"YurtMart/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Customer struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name"`
	Password        string             `json:"password"`
	ShippingAddress Address            `json:"shipping_address"`
	BillingAddress  Address            `json:"billing_address"`
	Admin           bool               `json:"admin" bson:"admin,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
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

func (c *Customer) GetShippingAddress() string {
	if c.ShippingAddress.Street == "" {
		return "No shipping address provided"
	}
	return c.ShippingAddress.Street + ", " +
		c.ShippingAddress.City + ", " +
		c.ShippingAddress.State + " " +
		c.ShippingAddress.PostalCode + ", " +
		c.ShippingAddress.Country
}
