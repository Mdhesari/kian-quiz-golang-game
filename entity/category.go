package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}
