package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Score float64

type UserScore struct {
	UserId      primitive.ObjectID `json:"user_id"`
	DisplayName string             `json:"display_name"`
	Score       Score              `json:"score"`
}
