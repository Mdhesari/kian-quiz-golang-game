package gameservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"time"

	"go.uber.org/zap"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/pkg/score"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
)

func (s *Service) GetAllGames(ctx context.Context, req param.GameGetAllRequest) (param.GameGetAllResponse, error) {
	op := "Game Service: Get all games."

	games, err := s.repo.GetAllGames(ctx, req.UserID)
	if err != nil {

		return param.GameGetAllResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GameGetAllResponse{
		Items: games,
	}, nil
}

func (s *Service) Create(ctx context.Context, req param.GameCreateRequest) (param.GameCreateResponse, error) {
	op := "Game Service: Create a new game."

	game := entity.Game{
		CategoryID: req.Category.ID,
		Questions:  req.Questions,
		Players:    req.Players,
		Status:     entity.GameStatusInProgress,
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

func (s *Service) GetGameById(ctx context.Context, req param.GameGetRequest) (param.GameGetResponse, error) {
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

func (s *Service) AnswerQuestion(ctx context.Context, req param.GameAnswerQuestionRequest) (param.GameAnswerQuestionResponse, error) {
	op := "Game service: answer question."

	gameRes, err := s.GetGameById(ctx, param.GameGetRequest{
		GameId: req.GameId,
	})
	if err != nil {
		logger.L().Error("Could not find the game.", zap.Error(err), zap.String("gameId", req.GameId.Hex()))

		return param.GameAnswerQuestionResponse{}, err
	}

	game := gameRes.Game
	if !game.IsInProgress() {

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrGameNotInProgress).WithKind(richerror.KindForbidden)
	}

	player, ok := game.Players[req.UserId.Hex()]
	if !ok {
		logger.L().Error("Player not found in the game.", zap.String("userId", req.UserId.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrGamePlayerNotFound).WithErr(err).WithKind(richerror.KindForbidden)
	}
	if !player.Status.InProgress() {

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrPlayerNotInProgress).WithErr(err).WithKind(richerror.KindForbidden)
	}

	playerAnswer := entity.PlayerAnswer{
		QuestionID: req.QuestionId,
		Answer: entity.Answer{
			Title: req.Answer,
		},
		EndTime: time.Now(),
	}

	if player.LastQuestionID != req.QuestionId {
		logger.L().Info("Player's last question ID does not match the question ID provided.", zap.String("questionId", playerAnswer.QuestionID.Hex()), zap.String("lastQuestionId", player.LastQuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrNotAnsweringCurrentQuestion).WithErr(err).WithKind(richerror.KindForbidden)
	}

	if player.HasAnsweredQuestion(playerAnswer.QuestionID) {
		logger.L().Info("Player has already answered this question.", zap.String("questionId", playerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrAlreadyAnswered).WithErr(err).WithKind(richerror.KindForbidden)
	}

	playerAnswer.StartTime = player.LastQuestionStartTime

	// TODO - Maybe better to handle this in repo
	var question entity.Question = game.GetQuestion(playerAnswer.QuestionID)
	if question.ID.IsZero() {
		logger.L().Error("Could not find the question.", zap.String("questionId", playerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, errmsg.ErrQuestionNotFound).WithKind(richerror.KindNotFound)
	}

	var correctAns entity.Answer = question.GetCorrectAnswer()
	if s.isCorrectAnswer(playerAnswer, correctAns) {
		playerAnswer.Answer.IsCorrect = true

		s := score.ClaculateScore(s.cfg.MaxScorePerQuestion, playerAnswer.GetAnswerTime(), s.cfg.MaxQuestionTimeout)
		playerAnswer.Score = s
	}

	ans, err := s.repo.CreateQuestionAnswer(ctx, req.UserId, req.GameId, playerAnswer)
	if err != nil {
		logger.L().Error(err.Error(), zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.String("question_id", playerAnswer.QuestionID.Hex()))

		return param.GameAnswerQuestionResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GameAnswerQuestionResponse{
		Answer:        ans,
		CorrectAnswer: correctAns,
	}, nil
}

func (s *Service) isCorrectAnswer(ans entity.PlayerAnswer, correctAns entity.Answer) bool {
	return correctAns.IsValid() && correctAns.Title == ans.Answer.Title
}

func (s *Service) GetNextQuestion(ctx context.Context, req param.GameGetNextQuestionRequest) (param.GameGetNextQuestionResponse, error) {
	op := "Game service: get next question."

	var nextQuestion entity.Question

	gameRes, err := s.GetGameById(ctx, param.GameGetRequest{
		GameId: req.GameId,
	})
	if err != nil {
		logger.L().Error("Could not find the game.", zap.Error(err), zap.String("gameId", req.GameId.Hex()))

		return param.GameGetNextQuestionResponse{}, richerror.New(op, errmsg.ErrGameNotFound).WithErr(err)
	}

	game := gameRes.Game
	if !game.IsInProgress() {

		return param.GameGetNextQuestionResponse{}, richerror.New(op, errmsg.ErrGameNotInProgress).WithKind(richerror.KindForbidden)
	}

	player, ok := game.Players[req.UserId.Hex()]
	if !ok {
		logger.L().Error("Player not found in the game.", zap.String("userId", req.UserId.Hex()))

		return param.GameGetNextQuestionResponse{}, richerror.New(op, errmsg.ErrGamePlayerNotFound).WithErr(err).WithKind(richerror.KindForbidden)
	}
	if !player.Status.InProgress() {

		return param.GameGetNextQuestionResponse{}, richerror.New(op, errmsg.ErrPlayerNotInProgress).WithErr(err).WithKind(richerror.KindForbidden)
	}

	for _, q := range game.Questions {
		if !player.HasAnsweredQuestion(q.ID) {
			nextQuestion = q

			break
		}
	}

	if nextQuestion.ID.IsZero() {
		logger.L().Info("No more questions available in the game.")

		return param.GameGetNextQuestionResponse{}, richerror.New(op, errmsg.ErrAllQuestionsAnswered).WithKind(richerror.KindOK)
	}

	player.LastQuestionID = nextQuestion.ID
	player.LastQuestionStartTime = time.Now()
	s.repo.UpdatePlayer(ctx, req.GameId, req.UserId, player)

	return param.GameGetNextQuestionResponse{
		Question: nextQuestion,
	}, nil
}

func (s *Service) FinishGame(ctx context.Context, req param.GameFinishRequest) (param.GameFinishResponse, error) {
	op := "Game service: finish game."

	if err := s.repo.UpdateGameStatus(ctx, req.GameId, entity.GameStatusFinished); err != nil {
		logger.L().Error(errmsg.ErrGameNotUpdated, zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.String("status", slice.GetGameStatusLabel(entity.GameStatusFinished)))

		return param.GameFinishResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	if err := s.repo.UpdateGameEndtime(ctx, req.GameId, time.Now()); err != nil {
		logger.L().Error(errmsg.ErrGameNotModified)

		return param.GameFinishResponse{}, err
	}

	payload := protobufencoder.EncodeGameFinishedEvent(entity.GameFinished{
		GameID: req.GameId,
	})
	s.pub.Publish(ctx, string(entity.GameStatusFinishedEvent), payload)

	return param.GameFinishResponse{}, nil
}

func (s *Service) UpdateWinner(ctx context.Context, req param.GameUpdateWinnerRequest) (param.GameUpdateWinnerResponse, error) {
	op := "Game service: update winner."

	if err := s.repo.UpdateGameWinner(ctx, req.GameId, req.Player); err != nil {
		logger.L().Error(errmsg.ErrGameNotUpdated, zap.Error(err), zap.String("game_id", req.GameId.Hex()), zap.Any("player", req.Player))

		return param.GameUpdateWinnerResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.GameUpdateWinnerResponse{}, nil
}
