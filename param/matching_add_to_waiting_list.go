package param

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchingAddToWaitingListRequest struct {
	UserID     primitive.ObjectID `json:"user_id" form:"user_id"`
	CategoryID primitive.ObjectID `json:"category_id" form:"category_id"`
}

type MatchingAddToWaitingListResponse struct {
	Timeout time.Duration `json:"timeout_in_nanoseconds"`
}
