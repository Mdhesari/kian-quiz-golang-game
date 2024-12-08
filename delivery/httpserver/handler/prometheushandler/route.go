package prometheushandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	r.GET("/metrics", h.Metrics)
}
