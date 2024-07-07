package matchinghandler

import (
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/claim"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) AddToWaitingList(c echo.Context) error {
	var req param.MatchingAddToWaitingListRequest
	if err := c.Bind(&req); err != nil {

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
