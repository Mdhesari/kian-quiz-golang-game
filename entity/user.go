package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Password []byte

type User struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	RoleID   *primitive.ObjectID `bson:"role_id,omitempty" json:"role_id"`
	Name     string              `bson:"name" json:"name"`
	Avatar   string              `bson:"avatar" json:"avatar"`
	Email    string              `bson:"email" json:"email"`
	Mobile   string              `bson:"mobile" json:"mobile"`
	Score    Score               `bson:"score" json:"score"`
	Password Password            `bson:"password" json:"-"`
}
