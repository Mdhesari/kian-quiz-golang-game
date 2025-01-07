package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answer struct {
	Title     string `bson:"title" json:"title"`
	IsCorrect bool   `bson:"is_correct" json:"-"`
}

type Question struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	CategoryID  primitive.ObjectID `bson:"category_id" json:"category_id"`
	Answers     []Answer           `bson:"answers" json:"answers"`
	Difficulty  QuestionDifficulty `bson:"difficulty" json:"difficulty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
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

func (q Question) GetCorrectAnswer() Answer {
	for _, a := range q.Answers {
		if a.IsCorrect {

			return a
		}
	}

	return Answer{}
}

func (a Answer) IsValid() bool {
	return a.Title != ""
}
