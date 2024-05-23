package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Answer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	QuestionID primitive.ObjectID
	Text       string
}
