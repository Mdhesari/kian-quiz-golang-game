package userhandler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Profile(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")

	claims, err := h.authSrv.VerifyToken(token)
	if err != nil {
		log.Println(err)

		return echo.ErrBadRequest
	}

	res, err := h.userSrv.GetByID(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	c.JSON(http.StatusOK, res)

	return nil
}
