package backpanelhandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	g := r.Group("/backpanel")

	g.GET("/users", h.ListUsers, middleware.Auth(h.authSrv, h.authConfig), middleware.RBAC(h.userSrv, h.rbacSrv))
}
