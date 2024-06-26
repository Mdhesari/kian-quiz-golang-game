package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Password []byte

type User struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty"`
	RoleID   *primitive.ObjectID `bson:"role_id,omitempty"`
	Name     string              `bson:"name"`
	Email    string              `bson:"email"`
	Mobile   string              `bson:"mobile"`
	Password Password            `bson:"password"`
}
