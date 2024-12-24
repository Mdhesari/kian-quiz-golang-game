package param

import "mdhesari/kian-quiz-golang-game/entity"

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	User  entity.User `json:"user"`
	Token string      `json:"token"`
}
