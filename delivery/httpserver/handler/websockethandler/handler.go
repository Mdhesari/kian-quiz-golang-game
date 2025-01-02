package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/pubsub"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
)

type Handler struct {
	pubsubManager *pubsub.PubSubManager
	presenceSrv   *presenceservice.Service
	authSrv       *authservice.Service
	authCfg       *authservice.Config
}

func New(pubsubMgr *pubsub.PubSubManager, presenceSrv *presenceservice.Service, authSrv *authservice.Service, authCfg *authservice.Config) Handler {
	return Handler{
		pubsubManager: pubsubMgr,
		presenceSrv:   presenceSrv,
		authSrv:       authSrv,
		authCfg:       authCfg,
	}
}
