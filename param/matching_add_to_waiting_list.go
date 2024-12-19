package param

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchingAddToWaitingListRequest struct {
	UserID     primitive.ObjectID `json:"user_id" form:"user_id"`
	CategoryID primitive.ObjectID `json:"category_id" form:"category_id"`
}

type MatchingAddToWaitingListResponse struct {
	Timeout uint `json:"timeout_in_seconds"`
}
