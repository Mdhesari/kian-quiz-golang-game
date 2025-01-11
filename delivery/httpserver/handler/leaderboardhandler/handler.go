package leaderboardhandler

import (
	"mdhesari/kian-quiz-golang-game/service/leaderboardservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
)

type Handler struct {
	leaderboardSrv *leaderboardservice.Service
	userSrv        *userservice.Service
}

func New(leaderboardSrv *leaderboardservice.Service, userSrv *userservice.Service) Handler {
	return Handler{
		leaderboardSrv: leaderboardSrv,
		userSrv:        userSrv,
	}
}
