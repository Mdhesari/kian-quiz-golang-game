package websockethandler

import (
	// "mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	group := r.Group("/ws")

	group.GET("", h.Websocket, middleware.Auth(h.authSrv, *h.authCfg))
}
