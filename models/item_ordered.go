package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemOrdered struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ItemID   primitive.ObjectID `bson:"item_id" json:"item_id"`
	Quantity int                `bson:"quantity" json:"quantity"`
	Price    float64            `bson:"price" json:"price"`
}

type Order struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Items []ItemOrdered      `bson:"items" json:"items"`
	Total float64            `bson:"total" json:"total"`
}
