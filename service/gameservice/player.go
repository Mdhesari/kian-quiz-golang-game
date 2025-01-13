package gameservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) UpdatePlayerStatus(ctx context.Context, req param.PlayerStatusUpdateRequest) (param.PlayerStatusUpdateResponse, error) {
	op := "Game Service: Update player status."

	modified, err := s.repo.UpdatePlayerStatus(ctx, req.GameId, req.UserId, req.Status)
	if err != nil {
		logger.L().Error("Could not update player status.", zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.String("user_id", req.UserId.Hex()))

		return param.PlayerStatusUpdateResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	if modified {
		logger.L().Info("Player status updated.", zap.String("game_id", req.GameId.Hex()), zap.String("user_id", req.UserId.Hex()), zap.String("status", req.Status.String()))

		if req.Status.Completed() {
			p := protobufencoder.EncodePlayerFinishedEvent(entity.PlayerFinished{
				UserId: req.UserId,
				GameId: req.GameId,
			})
			s.pub.Publish(ctx, string(entity.PlayerFinishedEvent), p)
		}
	}

	return param.PlayerStatusUpdateResponse{}, nil
}

func (s *Service) IncPlayerScore(ctx context.Context, req param.GamePlayerIncScoreRequest) (param.GamePlayerIncScoreResponse, error) {
	op := "Game service: increment player score."

	if err := s.repo.IncPlayerScore(ctx, req.GameId, req.UserId, req.Score); err != nil {
		logger.L().Error(errmsg.ErrGameNotUpdated, zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.String("user_id", req.UserId.Hex()), zap.Any("score", req.Score))

		return param.GamePlayerIncScoreResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GamePlayerIncScoreResponse{}, nil
}
