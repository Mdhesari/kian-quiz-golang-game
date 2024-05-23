package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}
