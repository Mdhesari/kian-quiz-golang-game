package matchinghandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	g := r.Group("/matching")

	g.POST("/add-to-waiting-list", h.AddToWaitingList, middleware.Auth(h.authSrv, h.authConfig), middleware.Presence(h.presenceSrv))
}
