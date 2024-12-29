package gameservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Repository interface {
	Create(ctx context.Context, game entity.Game) (entity.Game, error)
	GetGameById(ctx context.Context, id primitive.ObjectID) (entity.Game, error)
	Update(ctx context.Context, game entity.Game) error
	GetAllGames(ctx context.Context, UserID primitive.ObjectID) ([]entity.Game, error)
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

	return param.GameGetAllResponse{
		Items: games,
	}, nil
}

func (s Service) Create(ctx context.Context, req param.GameCreateRequest) (param.GameCreateResponse, error) {
	op := "Game Service: Create a new game."

	game := entity.Game{
		CategoryID: req.Category.ID,
		Questions:  req.Questions,
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