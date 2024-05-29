package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Products struct {
	Product_ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductName   string             `json:"product_name" validate:"required,max=255"`
	ProductPrice  float64            `json:"product_price" bson:"productprice" validate:"required,max=32"`
	Product_type  string             `json:"product_type" validate:"required,max=32"`
	Product_Point int                `json:"product_points" validate:"required,max=32"`
}
