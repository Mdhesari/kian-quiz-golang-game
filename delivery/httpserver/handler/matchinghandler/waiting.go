package matchinghandler

import (
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/claim"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h Handler) AddToWaitingList(c echo.Context) error {
	var req param.MatchingAddToWaitingListRequest
	if err := c.Bind(&req); err != nil {
		logger.L().Error("Matching handler error: %v\n", zap.Error(err))

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims := claim.GetClaimsFromEchoContext(c)
	req.UserID = claims.UserID

	if fields, err := h.matchingValidator.ValidateAddToWaitingListRequest(req); err != nil {
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fields,
		})
	}

	res, err := h.matchingSrv.AddToWaitingList(req)
	if err != nil {
		msg, code := richerror.Error(err)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, res)
}
