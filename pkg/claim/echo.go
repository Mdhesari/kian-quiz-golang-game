package claim

import (
	"mdhesari/kian-quiz-golang-game/pkg/constant"
	"mdhesari/kian-quiz-golang-game/service/authservice"

	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *authservice.Claims {
	return c.Get(constant.AuthContextKey).(*authservice.Claims)
}
