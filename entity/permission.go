package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string
}
