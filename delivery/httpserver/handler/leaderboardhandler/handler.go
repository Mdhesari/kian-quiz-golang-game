package leaderboardhandler

import "mdhesari/kian-quiz-golang-game/service/leaderboardservice"

type Handler struct {
	leaderboardSrv *leaderboardservice.Service
}

func New(leaderboardSrv *leaderboardservice.Service) Handler {
	return Handler{
		leaderboardSrv: leaderboardSrv,
	}
}
