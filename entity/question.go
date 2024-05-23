package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty"`
	Title           string               `bson:"title"`
	Description     string               `bson:"description"`
	CategoryID      primitive.ObjectID   `bson:"category_id"`
	AnswerIDs       []primitive.ObjectID `bson:"answer_ids"`
	CorrectAnswerID primitive.ObjectID   `bson:"correct_answer_id"`
	Difficulty      QuestionDifficulty   `bson:"difficulty"`
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
