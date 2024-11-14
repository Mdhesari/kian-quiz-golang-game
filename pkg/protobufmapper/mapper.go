package protobufmapper

import (
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/presence"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapFromProtobufPresenceResponseToParam(res *presence.GetPresenceResponse) param.PresenceResponse {
	paramRes := param.PresenceResponse{}

	for _, item := range res.Items {
		userId, _ := primitive.ObjectIDFromHex(item.UserId)
		paramRes.Items = append(paramRes.Items, param.PresenceItem{
			UserId:    userId,
			Timestamp: int64(item.Timestamp),
		})
	}

	return paramRes
}

func MapFromParamPresenceResponseToProtobuf(res param.PresenceResponse) *presence.GetPresenceResponse {
	protobufRes := presence.GetPresenceResponse{}

	for _, item := range res.Items {
		protobufRes.Items = append(protobufRes.Items, &presence.GetPresenceItem{
			UserId:    item.UserId.Hex(),
			Timestamp: uint64(item.Timestamp),
		})
	}

	return &protobufRes
}

func MapFromProtobufPresenceRequestToParam(req *presence.GetPresenceRequest) param.PresenceRequest {
	paramReq := param.PresenceRequest{}

	for _, userId := range req.UserId {
		pUserId, _ := primitive.ObjectIDFromHex(userId)
		paramReq.UserIds = append(paramReq.UserIds, pUserId)
	}

	return paramReq
}
