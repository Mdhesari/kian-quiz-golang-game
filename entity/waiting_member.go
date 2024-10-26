package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WaitingMember struct {
	UserId    primitive.ObjectID
	Category  Category
	Timestamp int64
}
