package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"mdhesari/kian-quiz-golang-game/websockethub"
)

type Handler struct {
	hub         *websockethub.Hub
	presenceSrv *presenceservice.Service
	authSrv     *authservice.Service
	authCfg     *authservice.Config
}

func New(hub *websockethub.Hub, presenceSrv *presenceservice.Service, authSrv *authservice.Service, authCfg *authservice.Config) Handler {
	return Handler{
		hub:         hub,
		presenceSrv: presenceSrv,
		authSrv:     authSrv,
		authCfg:     authCfg,
	}
}
