package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Access struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ActorID      string             `bson:"actor_id"`
	ActorType    string             `bson:"actory_type"`
	PermissionID string             `bson:"permission_id"`
}
