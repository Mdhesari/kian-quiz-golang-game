package param

import "mdhesari/kian-quiz-golang-game/entity"

type QuestionGetRequest struct {
	Count int `json:"count"`
}

type QuestionGetResponse struct {
	Items []entity.Question `json:"questions"`
}
