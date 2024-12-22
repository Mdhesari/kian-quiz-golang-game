package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"time"
)

type Config struct {
	PingPeriod   time.Duration `koan:"ping_period"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
}

type Handler struct {
	cfg       Config
	redisAdap *redisadapter.Adapter
	authSrv   *authservice.Service
	authCfg   *authservice.Config
}

func New(redisAdap *redisadapter.Adapter, authSrv *authservice.Service, authCfg *authservice.Config) Handler {
	return Handler{
		redisAdap: redisAdap,
		authSrv:   authSrv,
		authCfg:   authCfg,
	}
}
