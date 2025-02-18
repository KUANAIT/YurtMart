package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Payment struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	OrderID    primitive.ObjectID `json:"order_id,omitempty" bson:"order_id,omitempty"`
	Amount     float64            `json:"amount" bson:"amount"`
	Status     string             `json:"status" bson:"status"`
	Method     string             `json:"method" bson:"method"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}
