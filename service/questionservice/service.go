package questionservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo Repository
}

type Repository interface {
	GetRandomByCategory(ctx context.Context, categoryId primitive.ObjectID, count int) ([]entity.Question, error)
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) GetRandomQuestions(ctx context.Context, req param.QuestionGetRequest) (param.QuestionGetResponse, error) {
	op := "Question service: get random questions."

	items, err := s.repo.GetRandomByCategory(ctx, req.CategoryId, req.Count)
	if err != nil {

		return param.QuestionGetResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.QuestionGetResponse{
		Items: items,
	}, nil
}
