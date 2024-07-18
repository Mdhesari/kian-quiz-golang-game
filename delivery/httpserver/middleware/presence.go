package middleware

import (
	"log"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/claim"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Presence(presenceSrv presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			_, err = presenceSrv.Upsert(c.Request().Context(), param.PresenceUpsertRequest{
				UserId:    claims.UserID,
				Timestamp: timestamp.Now(),
			})
			if err != nil {
				log.Println("Presence Middleware Err: ", err)

				// There is a tradeoff here if we want to ignore presence errors or notice the user... we just notice for now
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrInternalServer,
				})
			}

			return next(c)
		}
	}
}