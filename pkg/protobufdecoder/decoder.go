package protobufdecoder

import (
	"encoding/base64"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/mongoutils"
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

		logger.L().Error("could not convert category hex to obj id {%v}.\n", zap.Error(err))
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

func DecodeGameStatusFinishedEvent(s string) entity.GameFinished {
	var pbE game.GameFinished
	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logger.L().Error("could not decode game status finished payload.", zap.Error(err))

		return entity.GameFinished{}
	}

	if err := protojson.Unmarshal(res, &pbE); err != nil {
		logger.L().Error("could not decode protobuf to json.", zap.Error(err))

		return entity.GameFinished{}
	}

	return entity.GameFinished{
		GameID: mongoutils.HexToObjectID(pbE.GameId),
	}
}

func DecodePlayerFinishedEvent(s string) entity.PlayerFinished {
	var pbE game.PlayerFinished

	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logger.L().Error("could not decode player finished payload.", zap.Error(err))

		return entity.PlayerFinished{}
	}

	if err := protojson.Unmarshal(res, &pbE); err != nil {
		logger.L().Error("could not decode protobuf to json.", zap.Error(err))

		return entity.PlayerFinished{}
	}

	return entity.PlayerFinished{
		UserId: mongoutils.HexToObjectID(pbE.UserId),
		GameId: mongoutils.HexToObjectID(pbE.GameId),
	}
}
