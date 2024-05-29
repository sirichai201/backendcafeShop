package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Orders struct {
	Order_id    primitive.ObjectID `json:"order_id" bson:"order_id"`
	User_id     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Product_id  primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Quantity    int64              `json:"quantity" bson:"quantity"`
	Total_Price float64            `json:"total_price" bson:"total_price"`
	Tel         string             `json:"tel" bson:"tel"`
}
