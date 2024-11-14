package param

import "mdhesari/kian-quiz-golang-game/entity"

type CategoryParam struct {
	//
}

type CategoryResponse struct {
	Items []entity.Category `json:"items"`
}
