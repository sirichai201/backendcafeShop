package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	
)

// type Points struct {
// 	User_id    primitive.ObjectID `json:"user_id"`
// 	Product_id primitive.ObjectID `json:"product_id"`
// 	Sum_point  int64              `json:"sum_point"`
// }
type Points struct {
	
	User_id    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Product_id primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Sum_point  int                `json:"sum_point,omitempty"`
	Phone      int                `json:"phone,omitempty"`
}
