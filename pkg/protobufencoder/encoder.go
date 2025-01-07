package protobufencoder

import (
	"encoding/base64"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/game"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/matching"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/websocket"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

func EncodePlayersMatchedEvent(e entity.PlayersMatched) string {
	pbE := matching.PlayersMatched{
		PlayerIds: slice.MapFromPrimitiveObjectIDToHexString(e.PlayerIDs),
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
		GameId: e.GameID.Hex(),
	}

	payload, err := protojson.Marshal(&pbE)
	if err != nil {
		// Update metrics
		logger.L().Error("Could not encode protobuf to json.")

		return ""
	}

	return base64.StdEncoding.EncodeToString(payload)
}

func EncodeWebSocketMsg(msg entity.WebsocketMsg) string {
	pbE := websocket.WebsocketMsg{
		Type:    msg.Type,
		Payload: msg.Payload,
	}

	res, err := protojson.Marshal(&pbE)
	if err != nil {
		logger.L().Error("Could not encoder ptobuf to json.", zap.Error(err))

		return ""
	}

	return base64.StdEncoding.EncodeToString(res)
}