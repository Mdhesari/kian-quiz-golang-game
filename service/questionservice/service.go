package questionservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
)

type Service struct {
	repo Repository
}

type Repository interface {
	GetRandomQuestions(ctx context.Context, count int) ([]entity.Question, error)
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) GetRandomQuestions(ctx context.Context, req param.QuestionGetRequest) (param.QuestionGetResponse, error) {
	op := "Question service: get random questions."

	items, err := s.repo.GetRandomQuestions(ctx, req.Count)
	if err != nil {

		return param.QuestionGetResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.QuestionGetResponse{
		Items: items,
	}, nil
}
