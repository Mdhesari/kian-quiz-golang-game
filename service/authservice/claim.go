package authservice

import (
	"mdhesari/kian-quiz-golang-game/pkg/constant"

	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) *Claims {
	return c.Get(constant.AuthContextKey).(*Claims)
}