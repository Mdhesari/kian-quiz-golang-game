package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/adapter/websocketadapter"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
)

type Handler struct {
	websocketAdapt *websocketadapter.Adapter
	presenceSrv    *presenceservice.Service
	authSrv        *authservice.Service
	authCfg        *authservice.Config
}

func New(websocketAdapter *websocketadapter.Adapter, presenceSrv *presenceservice.Service, authSrv *authservice.Service, authCfg *authservice.Config) Handler {
	return Handler{
		websocketAdapt: websocketAdapter,
		presenceSrv:    presenceSrv,
		authSrv:        authSrv,
		authCfg:        authCfg,
	}
}
