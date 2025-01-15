package categoryservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
)

type Service struct {
	repo Repistory
}

type Repistory interface {
	GetAll(ctx context.Context) ([]entity.Category, error)
}

func New(repo Repistory) Service {
	return Service{
		repo: repo,
	}
}

// TODO - Query param filters
func (s *Service) GetAll(ctx context.Context, _ param.CategoryParam) (param.CategoryResponse, error) {
	op := "Category service get:all"

	categories, err := s.repo.GetAll(ctx)
	if err != nil {

		return param.CategoryResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.CategoryResponse{
		Items: categories,
	}, nil
}