package gameservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

const (
	MaxQuestionTimeout  time.Duration = time.Second * 1500
	MaxScorePerQuestion uint8         = 5
)

type Config struct {
	GameTimeout time.Duration `koanf:"game_timeout"`
}

type Repository interface {
	Create(ctx context.Context, game entity.Game) (entity.Game, error)
	GetGameById(ctx context.Context, id primitive.ObjectID) (entity.Game, error)
	UpdatePlayer(ctx context.Context, gameId primitive.ObjectID, userId primitive.ObjectID, player entity.Player) error
	GetAllGames(ctx context.Context, userID primitive.ObjectID) ([]entity.Game, error)
	CreateQuestionAnswer(ctx context.Context, userId primitive.ObjectID, gameId primitive.ObjectID, playerAnswer entity.PlayerAnswer) (entity.PlayerAnswer, error)
	UpdateGameStatus(ctx context.Context, gameId primitive.ObjectID, status entity.GameStatus) error
	UpdateGameEndtime(ctx context.Context, gameId primitive.ObjectID, endTime time.Time) error
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
		Status:     entity.GameStatusInProgress,
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

	game := gameRes.Game
	if !game.IsInProgress() {

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrGameNotInProgress).WithKind(richerror.KindForbidden)
	}

	playerAnswer := entity.PlayerAnswer{
		QuestionID: req.QuestionId,
		Answer: entity.Answer{
			Title: req.Answer,
		},
		EndTime: time.Now(),
	}

	player, ok := game.Players[req.UserId.Hex()]
	if !ok {
		logger.L().Error("Player not found in the game.", zap.String("userId", req.UserId.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrGamePlayerNotFound).WithErr(err).WithKind(richerror.KindForbidden)
	}

	if player.HasAnsweredQuestion(playerAnswer.QuestionID) {
		logger.L().Info("Player has already answered this question.", zap.String("questionId", playerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrAlreadyAnswered).WithErr(err).WithKind(richerror.KindForbidden)
	}

	playerAnswer.StartTime = player.LastQuestionStartTime

	// TODO - Maybe better to handle this in repo
	var question entity.Question = game.GetQuestion(playerAnswer.QuestionID)
	if question.ID.IsZero() {
		logger.L().Error("Could not find the question.", zap.String("questionId", playerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrQuestionNotFound).WithKind(richerror.KindNotFound)
	}

	var correctAns entity.Answer = question.GetCorrectAnswer()
	if s.isCorrectAnswer(playerAnswer, correctAns) {
		playerAnswer.Answer.IsCorrect = true
		playerAnswer.Score = MaxScorePerQuestion
	}

	ans, err := s.repo.CreateQuestionAnswer(ctx, req.UserId, req.GameId, playerAnswer)
	if err != nil {
		logger.L().Error(err.Error(), zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.String("question_id", playerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GameAnswerQuestionResponse{
		Answer:        ans,
		CorrectAnswer: correctAns,
	}, nil
}

func (s *Service) isCorrectAnswer(ans entity.PlayerAnswer, correctAns entity.Answer) bool {
	return correctAns.IsValid() &&
		correctAns.Title == ans.Answer.Title &&
		!ans.IsTimeLimitReached(MaxQuestionTimeout)
}

func (s *Service) GetNextQuestion(ctx context.Context, req param.GameGetNextQuestionRequest) (param.GameGetNextQuestionResponse, error) {
	op := "Game service: get next question."

	var nextQuestion entity.Question

	gameRes, err := s.GetGameById(ctx, param.GameGetRequest{
		GameId: req.GameId,
	})
	if err != nil {
		logger.L().Error("Could not find the game.", zap.Error(err), zap.String("gameId", req.GameId.Hex()))

		return param.GameGetNextQuestionResponse{}, richerror.New(op, errmsg.ErrGameNotFound).WithErr(err)
	}

	game := gameRes.Game
	if !game.IsInProgress() {

		return param.GameGetNextQuestionResponse{}, richerror.New(op, errmsg.ErrGameNotInProgress).WithKind(richerror.KindForbidden)
	}

	player := game.Players[req.UserId.Hex()]
	for _, q := range game.Questions {
		if !player.HasAnsweredQuestion(q.ID) {
			nextQuestion = q

			break
		}
	}

	s.repo.UpdatePlayer(ctx, req.GameId, req.UserId, player)

	return param.GameGetNextQuestionResponse{
		Question: nextQuestion,
	}, nil
}

func (s *Service) FinishGame(ctx context.Context, req param.GameFinishRequest) (param.GameFinishResponse, error) {
	op := "Game service: finish game."

	if err := s.repo.UpdateGameStatus(ctx, req.GameId, entity.GameStatusFinished); err != nil {
		logger.L().Error(errmsg.ErrGameNotUpdated, zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.String("status", slice.GetGameStatusLabel(entity.GameStatusFinished)))

		return param.GameFinishResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	if err := s.repo.UpdateGameEndtime(ctx, req.GameId, time.Now()); err != nil {
		logger.L().Error(errmsg.ErrGameNotModified)

		return param.GameFinishResponse{}, err
	}

	return param.GameFinishResponse{}, nil
}
