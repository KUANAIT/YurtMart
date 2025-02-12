package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShoppingCart struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Items      []ItemOrdered      `bson:"items" json:"items"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
