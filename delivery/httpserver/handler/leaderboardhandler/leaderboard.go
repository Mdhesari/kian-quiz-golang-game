package leaderboardhandler

import (
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *Handler) GetLeaderboard(c echo.Context) error {
	var req param.LeaderboardRequest
	if err := c.Bind(&req); err != nil {
		logger.L().Error("Could not bind leaderboard request.", zap.Error(err))

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	res, err := h.leaderboardSrv.GetLeaderboard(c.Request().Context(), req)
	if err != nil {
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"Message": msg,
		})
	}

	return c.JSON(http.StatusOK, res)
}