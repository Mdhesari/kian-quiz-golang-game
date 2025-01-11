package param

import "mdhesari/kian-quiz-golang-game/entity"

type LeaderboardRequest struct {
	Limit int `json:"limit" param:"limit"`
}

type LeaderboardResponse struct {
	Leaderboard []entity.UserScore `json:"leaderboard"`
}
