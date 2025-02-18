package gamehandler

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *Handler) GetGame(c echo.Context) error {
	var req param.GameGetRequest
	if err := c.Bind(&req); err != nil {
		logger.L().Error("Could not bind game request.", zap.Error(err))

		return echo.NewHTTPError(http.StatusBadGateway)
	}

	res, err := h.gameSrv.GetGameById(c.Request().Context(), req)
	if err != nil {
		logger.L().Error("Could not get game by game id.", zap.Error(err), zap.Any("param", req))

		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"Message": msg,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllGames(c echo.Context) error {
	userId := authservice.GetClaims(c).UserID

	var req param.GameGetAllRequest
	req.UserID = userId

	res, err := h.gameSrv.GetAllGames(c.Request().Context(), req)
	if err != nil {
		logger.L().Error("Could not get all games.", zap.Error(err))

		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"Message": msg,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetNextQuestion(c echo.Context) error {
	userId := authservice.GetClaims(c).UserID

	var req param.GameGetNextQuestionRequest
	req.UserId = userId
	if err := c.Bind(&req); err != nil {
		logger.L().Error("Could not bind game request.", zap.Error(err))

		return echo.NewHTTPError(http.StatusBadGateway)
	}

	res, err := h.gameSrv.GetNextQuestion(c.Request().Context(), req)
	if err != nil {
		logger.L().Error("Could not get next question.", zap.Error(err), zap.Any("param", req))

		msg, code := richerror.Error(err)

		if errmsg.ErrAllQuestionsAnswered == msg {
			go h.gameSrv.UpdatePlayerStatus(context.Background(), param.PlayerStatusUpdateRequest{
				GameId: req.GameId,
				UserId: req.UserId,
				Status: entity.PlayerStatusCompleted,
			})
		}

		return c.JSON(code, echo.Map{
			"Message": msg,
		})
	}

	if res.Question.ID.IsZero() {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "No more questions available",
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AnswerQuestion(c echo.Context) error {
	userId := authservice.GetClaims(c).UserID

	var req param.GameAnswerQuestionRequest
	req.UserId = userId
	if err := c.Bind(&req); err != nil {
		logger.L().Error("Could not bind game request.", zap.Error(err))

		return echo.NewHTTPError(http.StatusBadGateway)
	}

	if fields, err := h.gameValidator.ValidateAnswerQuestion(req); err != nil {
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fields,
		})
	}

	res, err := h.gameSrv.AnswerQuestion(c.Request().Context(), req)
	if err != nil {
		logger.L().Error("Could not answer question.", zap.Error(err), zap.Any("param", req))

		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"Message": msg,
		})
	}

	h.gameSrv.IncPlayerScore(c.Request().Context(), param.GamePlayerIncScoreRequest{
		GameId: req.GameId,
		UserId: userId,
		Score:  res.Answer.Score,
	})
	// TODO - Error handling for inc score

	return c.JSON(http.StatusOK, res)
}
