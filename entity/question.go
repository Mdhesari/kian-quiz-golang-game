package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Answer struct {
	Title     string `bson:"title"`
	IsCorrect bool   `bson:"is_correct"`
}

type Question struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	CategoryID  primitive.ObjectID `bson:"category_id"`
	Answers     []Answer           `bson:"answers"`
	Difficulty  QuestionDifficulty `bson:"difficulty"`
}

type QuestionDifficulty int

const (
	QuestionDifficultyEasy QuestionDifficulty = iota
	QuestionDifficultyMedium
	QuestionDifficultyHard
)

func (q QuestionDifficulty) IsValid() bool {
	return q >= QuestionDifficultyEasy && q <= QuestionDifficultyHard
}
