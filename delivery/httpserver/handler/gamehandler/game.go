package gamehandler

import (
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h Handler) GetGames(c echo.Context) error {
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
