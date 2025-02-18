package httpserver

import (
	"fmt"
	"mdhesari/kian-quiz-golang-game/logger"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Handler interface {
	SetRoutes(r *echo.Echo)
}

type Config struct {
	Port int `koanf:"port"`
}

type Server struct {
	config   Config
	Router   *echo.Echo
	handlers []Handler
}

func New(c Config, h []Handler) Server {
	return Server{
		Router:   echo.New(),
		handlers: h,
		config:   c,
	}
}

func (s Server) Serve() {
	s.Router.Use(middleware.CORS())
	s.Router.Use(echoprometheus.NewMiddleware("echo"))
	s.Router.Use(middleware.RequestID())
	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var err string = ""
			if v.Error != nil {
				err = v.Error.Error()
			}
			logger.L().Named("http-server").Info("request",
				zap.String("request_id", v.RequestID),
				zap.String("host", v.Host),
				zap.String("content-length", v.ContentLength),
				zap.String("protocol", v.Protocol),
				zap.String("method", v.Method),
				zap.Duration("latency", v.Latency),
				zap.String("error", err),
				zap.String("remote_ip", v.RemoteIP),
				zap.Int64("response_size", v.ResponseSize),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
		HandleError:      true,
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
	}))
	s.Router.Use(middleware.Recover())

	// config handlers
	for _, h := range s.handlers {
		h.SetRoutes(s.Router)
	}

	// Start server
	address := fmt.Sprintf(":%d", s.config.Port)
	logger.L().Info(fmt.Sprintf("start echo server on %s\n", address))
	err := s.Router.Start(address)
	if err != nil {
		logger.L().Error("Http serve failed.", zap.Error(err))
	}
}
