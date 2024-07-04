package matchinghandler

import (
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
)

type Handler struct {
	authConfig  authservice.Config
	authSrv     *authservice.Service
	matchingSrv matchingservice.Service
}

func New(authConfig authservice.Config, authSrv *authservice.Service, matchingSrv matchingservice.Service) Handler {
	return Handler{
		authConfig:  authConfig,
		authSrv:     authSrv,
		matchingSrv: matchingSrv,
	}
}
