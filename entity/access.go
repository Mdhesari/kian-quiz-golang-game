package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Access struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	RoleID       string             `bson:"role_id,omitempty"`
	PermissionID string             `bson:"permission_id,omitempty"`
}
