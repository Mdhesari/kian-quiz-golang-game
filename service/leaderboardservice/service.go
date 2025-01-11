package leaderboardservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	"go.uber.org/zap"
)

type Repository interface {
	GetLeaderboard(ctx context.Context, limit int) ([]entity.UserScore, error)
	UpsertLeaderboardUserScore(ctx context.Context, us entity.UserScore) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) GetLeaderboard(ctx context.Context, req param.LeaderboardRequest) (param.LeaderboardResponse, error) {
	op := "Leaderboard service: Get leaderboard."

	leaderboard, err := s.repo.GetLeaderboard(ctx, req.Limit)
	if err != nil {

		return param.LeaderboardResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	logger.L().Info("Leaderboard retrieved.", zap.Any("leaderboard", leaderboard))

	return param.LeaderboardResponse{
		Leaderboard: leaderboard,
	}, nil
}

func (s *Service) UpsertLeaderboardUserScore(ctx context.Context, us entity.UserScore) error {
	op := "Leaderboard service: Upsert leaderboard user score."

	err := s.repo.UpsertLeaderboardUserScore(ctx, us)
	if err != nil {
		logger.L().Error("Could not upsert leaderboard user score.", zap.Error(err), zap.Any("us", us))

		return richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	logger.L().Info("Leaderboard user score upserted.", zap.Any("us", us))

	return nil
}
