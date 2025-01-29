package models

type ProductUser struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" bson:"name"`
	Category    string  `json:"category" bson:"category"`
	Price       float64 `json:"price" bson:"price"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Description string  `json:"description" bson:"description"`
}
