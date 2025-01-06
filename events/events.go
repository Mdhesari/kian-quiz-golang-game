package events

import (
	"context"
	"errors"
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

const (
	MaxQuestionTimeout  time.Duration = time.Second * 15
	MaxScorePerQuestion uint8         = 5
)

type EventManager struct {
	hub           *websockethub.Hub
	pubsubManager *pubsub.PubSubManager
	gameSrv       *gameservice.Service
	userSrv       *userservice.Service
	questionSrv   *questionservice.Service
}

func New(hub *websockethub.Hub, pubsubManager *pubsub.PubSubManager, gameSrv *gameservice.Service, userSrv *userservice.Service, questionSrv *questionservice.Service) EventManager {
	return EventManager{
		hub:           hub,
		pubsubManager: pubsubManager,
		gameSrv:       gameSrv,
		userSrv:       userSrv,
		questionSrv:   questionSrv,
	}
}

func (e EventManager) SubscribeEventHandlers() {
	logger.L().Info("Subscribing event listeners.")

	go e.pubsubManager.Subscribe(string(entity.PlayersMatchedEvent), e.HandlePlayersMatched)

	go e.pubsubManager.Subscribe(string(entity.GameStartedEvent), e.HandleHubGameStarted)

	go e.pubsubManager.Subscribe(string(entity.GamePlayerAnsweredEvent), e.HandleGamePlayerAnswered)

	logger.L().Info("event listeneres are subscribed.")
}

func (e EventManager) HandleGamePlayerAnswered(ctx context.Context, topic string, payload string) error {
	logger.L().Info("Handling game player answered event.")

	playerAnswered := protobufdecoder.DecodeGamePlayerAnsweredEvent(payload)
	logger.L().Info("event player answered recieved", zap.Any("playeransw", playerAnswered))

	playerAnswer := entity.PlayerAnswer{
		QuestionID: playerAnswered.QuestionID,
		Answer:     playerAnswered.Answer,
		StartTime:  time.Now().Add(-10 * time.Second),
		EndTime:    time.Now(),
	}

	gameRes, err := e.gameSrv.GetGameById(ctx, param.GameGetRequest{
		GameId: playerAnswered.GameID,
	})
	if err != nil {
		logger.L().Error("Could not find the game.", zap.Error(err), zap.String("gameId", playerAnswered.GameID.Hex()))

		return err
	}

	var question entity.Question
	for _, q := range gameRes.Game.Questions {
		if q.ID.Hex() == playerAnswered.QuestionID.Hex() {
			question = q
			break
		}
	}
	if question.ID.IsZero() {
		logger.L().Error("Could not find the question.", zap.String("questionId", playerAnswered.QuestionID.Hex()))

		return errors.New("question not found")
	}

	var correctAns entity.Answer
	for _, a := range question.Answers {
		if a.IsCorrect {
			correctAns = a
			break
		}
	}
	if correctAns.Title != "" && playerAnswer.Answer.Title != "" && correctAns.Title == playerAnswered.Answer.Title && playerAnswer.EndTime.Sub(playerAnswer.StartTime) <= MaxQuestionTimeout {
		logger.L().Info("Player answered correctly.", zap.String("questionId", playerAnswered.QuestionID.Hex()))

		playerAnswer.Answer.IsCorrect = true
		playerAnswer.Score = MaxScorePerQuestion
	}

	res, err := e.gameSrv.AnswerQuestion(ctx, param.GameAnswerQuestionRequest{
		UserId:       playerAnswered.UserID,
		GameId:       playerAnswered.GameID,
		PlayerAnswer: playerAnswer,
	})
	if err != nil {
		logger.L().Error("Could not answer question.", zap.Error(err), zap.Any("playeransw", playerAnswered))

		return err
	}

	logger.L().Info("Answer question successfully proceed.", zap.Any("res", res))

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
