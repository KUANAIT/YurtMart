package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   primitive.ObjectID `bson:"user_id,omitempty"`
	Username string             `bson:"username"`
	ItemID   primitive.ObjectID `bson:"item_id,omitempty"`
	Rating   int                `bson:"rating"`
	Text     string             `bson:"text"`
}
