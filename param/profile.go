package param

import "mdhesari/kian-quiz-golang-game/entity"

type ProfileResponse struct {
	User   *entity.User `json:"user"`
	Errors []string
}
