package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/adapter/websocketadapter"
	"mdhesari/kian-quiz-golang-game/service/authservice"
)

type Handler struct {
	websocketAdapt *websocketadapter.Adapter
	authSrv        *authservice.Service
	authCfg        *authservice.Config
}

func New(websocketAdapter *websocketadapter.Adapter, authSrv *authservice.Service, authCfg *authservice.Config) Handler {
	return Handler{
		websocketAdapt: websocketAdapter,
		authSrv:        authSrv,
		authCfg:        authCfg,
	}
}
