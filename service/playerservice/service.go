package playerservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Config struct {
	// Add any configuration parameters if needed
}

type Service struct {
	cfg  Config
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, player entity.Player) (entity.Player, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (entity.Player, error)
	Update(ctx context.Context, player entity.Player) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

func New(cfg Config, repo Repository) Service {
	return Service{
		cfg:  cfg,
		repo: repo,
	}
}

func (s Service) CreatePlayer(ctx context.Context, req param.PlayerCreateRequest) (param.PlayerCreateResponse, error) {
	op := "Player Service: Create player"

	player := entity.Player{
		UserID:    req.UserID,
		GameID:    req.GameID,
		Answers:   []entity.PlayerAnswer{},
		Score:     0,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.CreatedAt,
	}

	createdPlayer, err := s.repo.Create(ctx, player)
	if err != nil {
		logger.L().Error("Failed to create player", zap.Error(err), zap.Any("player", player))

		return param.PlayerCreateResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.PlayerCreateResponse{
		Player: createdPlayer,
	}, nil
}

func (s Service) GetPlayerByID(ctx context.Context, req param.PlayerGetRequest) (param.PlayerGetResponse, error) {
	op := "Player Service: Get player by ID"

	player, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		logger.L().Error("Failed to get player", zap.Error(err))

		return param.PlayerGetResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindNotFound)
	}

	return param.PlayerGetResponse{
		Player: player,
	}, nil
}

func (s Service) UpdatePlayer(ctx context.Context, req param.PlayerUpdateRequest) error {
	op := "Player Service: Update player"

	player, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		logger.L().Error("Failed to get player for update", zap.Error(err))

		return richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindNotFound)
	}

	// Update player fields
	player.Answers = req.Answers
	player.Score = req.Score
	player.UpdatedAt = req.UpdatedAt

	err = s.repo.Update(ctx, player)
	if err != nil {
		logger.L().Error("Failed to update player", zap.Error(err))
		return richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (s Service) DeletePlayer(ctx context.Context, req param.PlayerDeleteRequest) error {
	op := "Player Service: Delete player"

	err := s.repo.Delete(ctx, req.ID)
	if err != nil {
		logger.L().Error("Failed to delete player", zap.Error(err))
		return richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
