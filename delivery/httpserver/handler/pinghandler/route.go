package pinghandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	r.GET("/ping", h.Ping)
}