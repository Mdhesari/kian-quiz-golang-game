package leaderboardservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	"go.uber.org/zap"
)

type Repository interface {
	GetLeaderboard(ctx context.Context) ([]entity.User, error)
}

type Service struct {
	repo Repository
}

func New() *Service {
	return &Service{}
}

func (s *Service) GetLeaderboard(ctx context.Context) ([]entity.User, error) {
	op := "Leaderboard service: Get leaderboard."

	leaderboard, err := s.repo.GetLeaderboard(ctx)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	logger.L().Info("leaderboards.", zap.Any("l", leaderboard))

	return []entity.User{}, nil
}
