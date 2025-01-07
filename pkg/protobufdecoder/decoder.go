package protobufdecoder

import (
	"encoding/base64"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/game"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/matching"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/websocket"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

func DecodePlayersMatchedEvent(s string) entity.PlayersMatched {
	var pbE matching.PlayersMatched

	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// TOOD - update metrics

		panic(err)
	}
	if err := protojson.Unmarshal(res, &pbE); err != nil {
		// TODO - update metrics

		panic(err)
	}

	categoryId, err := primitive.ObjectIDFromHex(pbE.Category.ID)
	if err != nil {
		// TODO - update metrics

		log.Fatalf("could not convert category hex to obj id {%v}.\n", err)
	}

	return entity.PlayersMatched{
		PlayerIDs: slice.MapFromHexIDStringToPrimitiveObject(pbE.PlayerIds),
		Category: entity.Category{
			ID:          categoryId,
			Title:       pbE.Category.Ttile,
			Description: pbE.Category.Description,
		},
	}
}

func DecodeGameStartedEvent(s string) entity.GameStarted {
	var pbE game.GameStarted

	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logger.L().Error("could not decode game started payload.", zap.Error(err))
	}

	if err := protojson.Unmarshal(res, &pbE); err != nil {
		panic(err)
	}

	gameId, err := primitive.ObjectIDFromHex(pbE.GameId)
	if err != nil {
		logger.L().Error("decode game started event: Could not get objwct id.", zap.Error(err))
	}

	return entity.GameStarted{
		GameID: gameId,
	}
}

func DecodeWebSocketMsg(s string) entity.WebsocketMsg {
	var pbE websocket.WebsocketMsg

	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logger.L().Error("could not decode websocket msg payload.", zap.Error(err))

		return entity.WebsocketMsg{}
	}

	if err := protojson.Unmarshal(res, &pbE); err != nil {
		logger.L().Error("could not decode protobuf to json.", zap.Error(err))

		return entity.WebsocketMsg{}
	}

	return entity.WebsocketMsg{
		Type:    pbE.Type,
		Payload: pbE.Payload,
	}
}
