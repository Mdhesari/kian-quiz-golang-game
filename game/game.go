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
	"time"

	"go.uber.org/zap"
)

type Game struct {
	pubsubManager *pubsub.PubSubManager
	gameSrv       *gameservice.Service
	userSrv       *userservice.Service
	questionSrv   *questionservice.Service
}

func New(pubsubManager *pubsub.PubSubManager, gameSrv *gameservice.Service, userSrv *userservice.Service, questionSrv *questionservice.Service) Game {
	return Game{
		pubsubManager: pubsubManager,
		gameSrv:       gameSrv,
		userSrv:       userSrv,
		questionSrv:   questionSrv,
	}
}

func (m Game) SubscribeEventHandlers() {
	go m.pubsubManager.Subscribe(string(entity.PlayersMatchedEvent), m.HandlePlayersMatched)
}

func (m Game) HandlePlayersMatched(ctx context.Context, topic string, payload string) error {
	playersMatched := protobufdecoder.DecodePlayersMatchedEvent(payload)

	questionRes, err := m.questionSrv.GetRandomQuestions(ctx, param.QuestionGetRequest{
		CategoryId: playersMatched.Category.ID,
	})
	if err != nil {

		logger.L().Error("Could not get random questions for creating game.", zap.Error(err), zap.Any("Event", playersMatched))
	}

	var players []entity.Player
	for _, id := range playersMatched.PlayerIDs {
		res, err := m.userSrv.GetByID(id)
		if err != nil {
			logger.L().Error("Could not get user for creating game.", zap.Error(err), zap.Any("userID", id))

			continue
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
