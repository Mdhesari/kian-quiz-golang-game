package gameservice

import (
	"context"
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

const (
	MaxQuestionTimeout  time.Duration = time.Second * 15
	MaxScorePerQuestion uint8         = 5
)

type Repository interface {
	Create(ctx context.Context, game entity.Game) (entity.Game, error)
	GetGameById(ctx context.Context, id primitive.ObjectID) (entity.Game, error)
	Update(ctx context.Context, game entity.Game) error
	GetAllGames(ctx context.Context, userID primitive.ObjectID) ([]entity.Game, error)
	CreateQuestionAnswer(ctx context.Context, userId primitive.ObjectID, gameId primitive.ObjectID, playerAnswer entity.PlayerAnswer) (entity.PlayerAnswer, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) GetAllGames(ctx context.Context, req param.GameGetAllRequest) (param.GameGetAllResponse, error) {
	op := "Game Service: Get all games."

	games, err := s.repo.GetAllGames(ctx, req.UserID)
	if err != nil {

		return param.GameGetAllResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	if games == nil {
		games = []entity.Game{}
	}

	return param.GameGetAllResponse{
		Items: games,
	}, nil
}

func (s Service) Create(ctx context.Context, req param.GameCreateRequest) (param.GameCreateResponse, error) {
	op := "Game Service: Create a new game."

	game := entity.Game{
		CategoryID: req.Category.ID,
		Questions:  req.Questions,
		Players:    req.Players,
		StartTime:  time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	game, err := s.repo.Create(ctx, game)
	if err != nil {

		return param.GameCreateResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GameCreateResponse{
		Game: game,
	}, nil
}

func (s Service) GetGameById(ctx context.Context, req param.GameGetRequest) (param.GameGetResponse, error) {
	op := "Game service: find game by id."

	game, err := s.repo.GetGameById(ctx, req.GameId)
	if err != nil {
		logger.L().Error("Could not get game by id.", zap.Error(err), zap.String("game_id", req.GameId.Hex()))

		return param.GameGetResponse{}, err
	}
	if game.ID.IsZero() {
		logger.L().Info("Game does not exists!", zap.String("game_id", req.GameId.Hex()))

		return param.GameGetResponse{}, richerror.New(op, errmsg.ErrNotFound).WithErr(err).WithKind(richerror.KindNotFound)
	}

	return param.GameGetResponse{
		Game: game,
	}, nil
}

func (s Service) AnswerQuestion(ctx context.Context, req param.GameAnswerQuestionRequest) (param.GameAnswerQuestionResponse, error) {
	op := "Game service: answer question."

	gameRes, err := s.GetGameById(ctx, param.GameGetRequest{
		GameId: req.GameId,
	})
	if err != nil {
		logger.L().Error("Could not find the game.", zap.Error(err), zap.String("gameId", req.GameId.Hex()))

		return param.GameAnswerQuestionResponse{}, err
	}

	req.PlayerAnswer.StartTime = gameRes.Game.Players[req.UserId.Hex()].LastQuestionStartTime

	var question entity.Question
	for _, q := range gameRes.Game.Questions {
		if q.ID.Hex() == req.PlayerAnswer.QuestionID.Hex() {
			question = q
			break
		}
	}
	if question.ID.IsZero() {
		logger.L().Error("Could not find the question.", zap.String("questionId", req.PlayerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, errors.New("question not found")
	}

	var correctAns entity.Answer
	for _, a := range question.Answers {
		if a.IsCorrect {
			correctAns = a
			break
		}
	}
	if correctAns.Title != "" && req.PlayerAnswer.Answer.Title != "" && correctAns.Title == req.PlayerAnswer.Answer.Title && req.PlayerAnswer.EndTime.Sub(req.PlayerAnswer.StartTime) <= MaxQuestionTimeout {
		logger.L().Info("Player answered correctly.", zap.String("questionId", req.PlayerAnswer.QuestionID.Hex()))

		req.PlayerAnswer.Answer.IsCorrect = true
		req.PlayerAnswer.Score = MaxScorePerQuestion
	}

	playerAw, err := s.repo.CreateQuestionAnswer(ctx, req.UserId, req.GameId, req.PlayerAnswer)
	if err != nil {
		logger.L().Error(err.Error(), zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.String("question_id", req.PlayerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GameAnswerQuestionResponse{
		PlayerAnswer: playerAw,
	}, nil
}
