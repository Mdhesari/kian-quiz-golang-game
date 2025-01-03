package game

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
	"mdhesari/kian-quiz-golang-game/pubsub"
	"mdhesari/kian-quiz-golang-game/service/gameservice"
	"mdhesari/kian-quiz-golang-game/service/questionservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"mdhesari/kian-quiz-golang-game/websockethub"
	"time"

	"go.uber.org/zap"
)

type Game struct {
	hub           *websockethub.Hub
	pubsubManager *pubsub.PubSubManager
	gameSrv       *gameservice.Service
	userSrv       *userservice.Service
	questionSrv   *questionservice.Service
}

func New(hub *websockethub.Hub, pubsubManager *pubsub.PubSubManager, gameSrv *gameservice.Service, userSrv *userservice.Service, questionSrv *questionservice.Service) Game {
	return Game{
		hub:           hub,
		pubsubManager: pubsubManager,
		gameSrv:       gameSrv,
		userSrv:       userSrv,
		questionSrv:   questionSrv,
	}
}

func (m Game) SubscribeEventHandlers() {
	go m.pubsubManager.Subscribe(string(entity.PlayersMatchedEvent), m.HandlePlayersMatched)

	go m.pubsubManager.Subscribe(string(entity.GameStartedEvent), m.HandleHubGameStarted)
}

func (m Game) HandleHubGameStarted(ctx context.Context, topic string, payload string) error {
	gse := protobufdecoder.DecodeGameStartedEvent(payload)
	game, err := m.gameSrv.GetGameById(ctx, param.GameGetRequest{
		GameId: gse.GameID,
	})
	if err != nil {
		logger.L().Error("Handler game started: Could not get game.", zap.Error(err), zap.String("gameID", gse.GameID.Hex()))

		return err
	}

	for _, player := range game.Game.Players {
		m.hub.BroadcastMessage(&websockethub.Message{
			Type:   topic,
			UserID: player.UserID.Hex(),
			Body:   []byte(payload),
		})
	}

	return nil
}

func (m Game) HandlePlayersMatched(ctx context.Context, topic string, payload string) error {
	playersMatched := protobufdecoder.DecodePlayersMatchedEvent(payload)

	questionRes, err := m.questionSrv.GetRandomQuestions(ctx, param.QuestionGetRequest{
		CategoryId: playersMatched.Category.ID,
	})
	if err != nil {
		logger.L().Error("Could not get random questions for creating game.", zap.Error(err), zap.Any("Event", playersMatched))

		return err
	}

	var players []entity.Player
	for _, id := range playersMatched.PlayerIDs {
		res, err := m.userSrv.GetByID(id)
		if err != nil {
			logger.L().Error("Could not get user for creating game.", zap.Error(err), zap.Any("userID", id))

			return err
		}

		players = append(players, entity.Player{
			Name:      res.User.Name,
			UserID:    res.User.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	// TODO - game questions should be nullable so that if in the starting game scenario we couldn't fetch questions finlally raise error to user.
	game, err := m.gameSrv.Create(context.Background(), param.GameCreateRequest{
		Category:  playersMatched.Category,
		Questions: questionRes.Items,
		Players:   players,
	})
	if err != nil {
		logger.L().Error("Could not setup a game.", zap.Error(err), zap.Any("game", game))

		return err
	}

	logger.L().Info("A new game created and updated with player IDs.", zap.Any("game", game.Game.ID))

	gameStartedPayload := protobufencoder.EncodeGameStartedEvent(entity.GameStarted{
		GameID: game.Game.ID,
	})

	m.pubsubManager.Publish(context.Background(), string(entity.GameStartedEvent), gameStartedPayload)

	return nil
}
