package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	SetRoutes(r *echo.Echo)
}

type HTTPServer struct {
	Port    int    `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer `koanf:"http_server"`
}

type Server struct {
	config   Config
	handlers []Handler
	Router   *echo.Echo
}

func New(c Config, r *echo.Echo, h []Handler) Server {
	return Server{
		handlers: h,
		config: c,
		Router: r,
	}
}

func (s Server) Serve() {
	// config handlers
	for _, h := range s.handlers {
		h.SetRoutes(s.Router)
	}

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)

	fmt.Printf("start echo server on %s\n", address)

	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}
