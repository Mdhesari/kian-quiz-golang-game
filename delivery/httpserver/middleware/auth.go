package middleware

import (
	"mdhesari/kian-quiz-golang-game/pkg/constant"
	"mdhesari/kian-quiz-golang-game/service/authservice"

	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey:    constant.AuthContextKey,
		SigningMethod: constant.AuthSigningMethod,
		SigningKey:    config.Secret,
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.VerifyToken(auth)
			if err != nil {

				return nil, err
			}

			return claims, nil
		},
	})
}
