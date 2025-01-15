package questionservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	"go.uber.org/zap"
)


func (s *Service) GetRandomQuestions(ctx context.Context, req param.QuestionGetRequest) (param.QuestionGetResponse, error) {
	op := "Question service: get random questions."

	logger.L().Info("Getting random questions.", zap.Int("count", s.cfg.QuestionsCount))

	items, err := s.repo.GetRandomByCategory(ctx, req.CategoryId, s.cfg.QuestionsCount)
	if err != nil {

		return param.QuestionGetResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.QuestionGetResponse{
		Items: items,
	}, nil
}
