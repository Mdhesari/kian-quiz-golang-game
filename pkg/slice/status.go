package slice

import "mdhesari/kian-quiz-golang-game/entity"

func GetGameStatusLabel(status entity.GameStatus) string {
	return []string{"Aborted", "In Progress", "Finished"}[status]
}
