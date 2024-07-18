package param

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PresenceUpsertRequest struct {
	UserId primitive.ObjectID
	Timestamp int64
}

type PrescenceUpsertResponse struct {
	//
}