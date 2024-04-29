package userhandler

import (
	"mdhesari/kian-quiz-golang-game/pkg/constant"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getClaims(c echo.Context) *authservice.Claims {
	return c.Get(constant.AuthContextKey).(*authservice.Claims)
}

func (h Handler) Profile(c echo.Context) error {
	claims := getClaims(c)

	res, err := h.userSrv.GetByID(claims.UserID)
	if err != nil {
		msg, code := richerror.Error(err)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, res)
}
