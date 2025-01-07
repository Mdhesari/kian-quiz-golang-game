package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerAnswer struct {
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
	Answer     Answer             `bson:"answer" json:"answer"`
	Score      uint8              `bson:"score" json:"score"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	EndTime    time.Time          `bson:"end_time" json:"end_time"`
}

type Player struct {
	Name                  string             `bson:"name" json:"name"`
	Answers               []PlayerAnswer     `bson:"answers" json:"answers"`
	Score                 uint8              `bson:"score" json:"score"`
	IsWinner              bool               `bson:"is_winner,omitempty" json:"is_winner,omitempty"`
	LastQuestionID        primitive.ObjectID `bson:"last_question_id,omitempty" json:"last_question_id,omitempty"`
	LastQuestionStartTime time.Time          `bson:"last_question_start_time,omitempty" json:"last_question_start_time,omitempty"`
	CreatedAt             time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time          `bson:"updated_at" json:"updated_at"`
}

type PlayerAnswered struct {
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	GameID     primitive.ObjectID `bson:"game_id" json:"game_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
	Answer     Answer             `bson:"answer" json:"answer"`
}

func (p *Player) HasAnsweredQuestion(questionID primitive.ObjectID) bool {
	for _, answer := range p.Answers {
		if answer.QuestionID == questionID {
			return true
		}
	}

	return false
}
