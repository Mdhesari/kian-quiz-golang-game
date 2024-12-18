package param

import "mdhesari/kian-quiz-golang-game/entity"

type CategoryParam struct {
	Title string `json:"title"`
}

type CategoryResponse struct {
	Items []entity.Category `json:"items"`
}
