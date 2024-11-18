package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlayersMatched struct {
	Players []primitive.ObjectID
	Category Category
}
