package param

import "go.mongodb.org/mongo-driver/bson/primitive"

type PresenceItem struct {
	UserId    primitive.ObjectID
	Timestamp int64
}

type PresenceResponse struct {
	Items []PresenceItem
}

type PresenceRequest struct {
	UserIds []primitive.ObjectID
}
