package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WaitingMember struct {
	UserId    primitive.ObjectID `json:"user_id`
	Category  Category           `json:"category"`
	Timestamp int64              `json:"timestamp"`
}
