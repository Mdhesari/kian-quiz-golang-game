package gameservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
)

type Repository interface {
	Create(ctx context.Context, game entity.Game) (entity.Game, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) Create(ctx context.Context, req param.GameCreateRequest) (param.GameCreateResponse, error) {
	op := "Game Service: Create a new game."

	game, err := s.repo.Create(ctx, entity.Game{
		PlayerIDs:  req.Players,
		CategoryID: req.Category.ID,
	})
	if err != nil {

		return param.GameCreateResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GameCreateResponse{
		Game: game,
	}, nil
}
