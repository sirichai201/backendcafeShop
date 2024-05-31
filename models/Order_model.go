package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Orders struct {
	Order_id    primitive.ObjectID `json:"order_id" bson:"order_id"`
	CreatedAt   time.Time          `json:"created_at"`
	Status      string             `json:"status" bson:"status"`
	User_id     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Product_id  primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Total_Price int                `json:"total_price" bson:"total_price"`
	Member      string             `json:"Member" bson:"Member"`
	PesentPoint int                `json:"persent_point"`
	
}

type Bill struct {
	User_id     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	Total_Price int                `json:"total_price" bson:"total_price"`
	Payment      string             `json:"payment"`
	PesentPoint int                `json:"persent_point"`
}
