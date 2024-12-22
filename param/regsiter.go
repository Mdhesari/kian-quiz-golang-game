package param

import "mdhesari/kian-quiz-golang-game/entity"

type RegisterRequest struct {
	Name     string `json:"name" form:"name"`
	Mobile   string `json:"mobile" form:"mobile"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterResponse struct {
	User *entity.User `json:"user"`
}
