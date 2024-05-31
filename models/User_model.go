package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	User_ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username" validate:"required,min=3,max=32"`
	Email    string             `json:"email" validate:"required,email,min=3,max=32"`
	Password string             `json:"password" validate:"required,min=3,max=32"`
	Role     string             `json:"role"`
	Status   int                `json:"status"`
	Phone    string             `json:"phone" validate:"required,min=8,max=10"`
	Point    int                `json:"point"`
}
