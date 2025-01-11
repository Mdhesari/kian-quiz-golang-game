package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Score uint

type UserScore struct {
	UserId primitive.ObjectID `json:"user_id"`
	Score  Score              `json:"score"`
}
