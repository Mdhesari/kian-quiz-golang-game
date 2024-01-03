package userhandler

import (
	"log"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Profile(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")

	claims, err := h.authSrv.VerifyToken(token)
	if err != nil {
		log.Println(err)

		return echo.NewHTTPError(http.StatusUnauthorized, param.ProfileResponse{
			User:   nil,
			Errors: []string{"Invalid Authorization Token."},
		})
	}

	res, err := h.userSrv.GetByID(claims.UserID)
	if err != nil {
		if err == userservice.ErrNotFound {

			return echo.ErrNotFound
		}

		return echo.NewHTTPError(500, err)
	}

	return c.JSON(http.StatusOK, res)
}
