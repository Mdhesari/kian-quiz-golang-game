package events

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/mongoutils"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
	"mdhesari/kian-quiz-golang-game/pubsub"
	"mdhesari/kian-quiz-golang-game/service/gameservice"
	"mdhesari/kian-quiz-golang-game/service/leaderboardservice"
	"mdhesari/kian-quiz-golang-game/service/questionservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"mdhesari/kian-quiz-golang-game/websockethub"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type EventManager struct {
	hub            *websockethub.Hub
	pubsubManager  *pubsub.PubSubManager
	gameSrv        *gameservice.Service
	userSrv        *userservice.Service
	questionSrv    *questionservice.Service
	leaderboardSrv *leaderboardservice.Service
	gameCfg        *gameservice.Config
}

func New(hub *websockethub.Hub, pubsubManager *pubsub.PubSubManager, gameSrv *gameservice.Service, userSrv *userservice.Service, questionSrv *questionservice.Service, leaderboardSrv *leaderboardservice.Service, gameCfg *gameservice.Config) EventManager {
	return EventManager{
		hub:            hub,
		pubsubManager:  pubsubManager,
		gameSrv:        gameSrv,
		userSrv:        userSrv,
		questionSrv:    questionSrv,
		leaderboardSrv: leaderboardSrv,
		gameCfg:        gameCfg,
	}
}

func (e EventManager) SubscribeEventHandlers() {
	logger.L().Info("Subscribing event listeners.")

	go e.pubsubManager.Subscribe(string(entity.PlayersMatchedEvent), e.HandlePlayersMatched)

	go e.pubsubManager.Subscribe(string(entity.GameStartedEvent), e.HandleHubGameStarted)

	go e.pubsubManager.Subscribe(string(entity.GameStatusFinishedEvent), e.HandleHubGameFinished)

	logger.L().Info("event listeneres are subscribed.")
}

func (e EventManager) HandleHubGameFinished(ctx context.Context, topic string, payload string) error {
	logger.L().Info("game status finished.", zap.String("topic", topic), zap.String("payload", payload))

	gse := protobufdecoder.DecodeGameStatusFinishedEvent(payload)
	gameRes, err := e.gameSrv.GetGameById(ctx, param.GameGetRequest{
		GameId: gse.GameID,
	})
	if err != nil {
		logger.L().Error("Handler game finished: Could not get game.", zap.Error(err), zap.String("gameID", gse.GameID.Hex()))

		return err
	}

	players := gameRes.Game.Players
	var winner entity.Player = entity.Player{}
	for userId, player := range players {
		err := e.userSrv.IncScore(ctx, param.UserIncrementScoreRequest{
			UserId: mongoutils.HexToObjectID(userId),
			Score:  player.Score,
		})
		if err != nil {
			logger.L().Error("Handler game finished: Could not increment score.", zap.Error(err), zap.String("userId", userId), zap.Any("player", player))
		}

		if player.Score > 0 {
			us := entity.UserScore{
				UserId:      mongoutils.HexToObjectID(userId),
				Score:       player.Score,
				DisplayName: player.Name,
			}
			if err := e.leaderboardSrv.UpsertLeaderboardUserScore(ctx, us); err != nil {
				logger.L().Error("Handler game finished: Could not upsert leaderboard.", zap.Error(err), zap.String("userId", userId), zap.Any("player", player))
			}
		}

		if winner.Score < player.Score {
			winner = player
		}
	}

	e.gameSrv.UpdateWinner(ctx, param.GameUpdateWinnerRequest{
		GameId: gse.GameID,
		Player: winner,
	})

	return nil
}

func (e EventManager) HandleHubGameStarted(ctx context.Context, topic string, payload string) error {
	gse := protobufdecoder.DecodeGameStartedEvent(payload)
	gameRes, err := e.gameSrv.GetGameById(ctx, param.GameGetRequest{
		GameId: gse.GameID,
	})
	if err != nil {
		logger.L().Error("Handler game started: Could not get game.", zap.Error(err), zap.String("gameID", gse.GameID.Hex()))

		return err
	}

	// TODO - Temporary - need queue or just leave it to other processes
	go func(gameSrv *gameservice.Service, gameId primitive.ObjectID, d time.Duration) {
		time.Sleep(d)

		logger.L().Info("Finishing the game...")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		gameSrv.FinishGame(ctx, param.GameFinishRequest{
			GameId: gameId,
		})

	}(e.gameSrv, gse.GameID, e.gameCfg.GameTimeout)

	logger.L().Info("Decoded game started event.", zap.Any("game", payload))

	var userIDs []string
	for userId := range gameRes.Game.Players {
		userIDs = append(userIDs, userId)
	}

	e.hub.BroadcastMessage(&websockethub.Message{
		Type:    topic,
		UserIDs: userIDs,
		Body:    payload,
	})

	logger.L().Info("Broadcasted message", zap.String("topic", topic))

	return nil
}

func (m EventManager) HandlePlayersMatched(ctx context.Context, topic string, payload string) error {
	playersMatched := protobufdecoder.DecodePlayersMatchedEvent(payload)

	questionRes, err := m.questionSrv.GetRandomQuestions(ctx, param.QuestionGetRequest{
		CategoryId: playersMatched.Category.ID,
	})
	if err != nil {
		logger.L().Error("Could not get random questions for creating game.", zap.Error(err), zap.Any("Event", playersMatched))

		return err
	}

	var players map[string]entity.Player = make(map[string]entity.Player)
	for _, id := range playersMatched.PlayerIDs {
		res, err := m.userSrv.GetByID(id)
		if err != nil {
			logger.L().Error("Could not get user for creating game.", zap.Error(err), zap.Any("userID", id))

			return err
		}

		players[res.User.ID.Hex()] = entity.Player{
			Name:      res.User.Name,
			Answers:   []entity.PlayerAnswer{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Score:     0,
			Status:    entity.PlayerStatusInProgress,
		}
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
