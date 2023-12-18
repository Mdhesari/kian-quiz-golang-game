package entity

type Question struct {
	ID              uint
	Title           string
	Description     string
	CategoryID      uint
	AnswerIDs       []uint
	CorrectAnswerID uint
	Difficulty      QuestionDifficulty
}

type QuestionDifficulty uint

const (
	QuestionDifficultyEasy QuestionDifficulty = iota
	QuestionDifficultyMedium
	QuestionDifficultyHard
)

func (q QuestionDifficulty) IsValid() bool {
	return q >= QuestionDifficultyEasy && q <= QuestionDifficultyHard
}
