package protobufencoder

import (
	"encoding/base64"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/game"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/matching"

	"google.golang.org/protobuf/encoding/protojson"
)

func EncodeUsersMatchedEvent(e entity.PlayersMatched) string {
	pbE := matching.UsersMatched{
		UserIds: slice.MapFromPrimitiveObjectIDToHexString(e.Players),
		Category: &matching.Category{
			ID:          e.Category.ID.Hex(),
			Ttile:       e.Category.Title,
			Description: e.Category.Description,
		},
	}

	payload, err := protojson.Marshal(&pbE)
	if err != nil {
		// TODO - update metrics
		// TODO - log error

		return ""
	}

	return base64.StdEncoding.EncodeToString(payload)
}

func EncodeGameStartedEvent(e entity.GameStarted) string {
	pbE := game.GameStarted{
		PlayerIds: slice.MapFromPrimitiveObjectIDToHexString(e.PlayerIds),
	}

	payload, err := protojson.Marshal(&pbE)
	if err != nil {
		// Update metrics
		logger.L().Error("Could not encode protobuf to json.")

		return ""
	}

	return base64.StdEncoding.EncodeToString(payload)
}
