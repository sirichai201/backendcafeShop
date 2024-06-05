package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Products struct {
	ProductID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Image        string             `json:"image" validate:"required`
	ProductName  string             `json:"product_name" validate:"required,max=255"`
	ProductPrice float64            `json:"product_price" bson:"productprice" validate:"required,max=32"`
	ProductType  string             `json:"product_type" validate:"required,max=32"`
	ProductPoint int                `json:"product_points" validate:"required,max=32"`
	Description string             `json:"description" validate:"required`
}
