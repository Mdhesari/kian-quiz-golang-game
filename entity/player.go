package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerStatus uint8

const (
	PlayerStatusRejected PlayerStatus = iota
	PlayerStatusInProgress
	PlayerStatusCompleted
)

type PlayerAnswer struct {
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
	Answer     Answer             `bson:"answer" json:"answer"`
	Score      Score              `bson:"score" json:"score"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	EndTime    time.Time          `bson:"end_time" json:"end_time"`
}

type Player struct {
	Name                  string             `bson:"name" json:"name"`
	Answers               []PlayerAnswer     `bson:"answers" json:"answers"`
	Score                 Score              `bson:"score" json:"score"`
	IsWinner              bool               `bson:"is_winner,omitempty" json:"is_winner,omitempty"`
	LastQuestionID        primitive.ObjectID `bson:"last_question_id,omitempty" json:"last_question_id,omitempty"`
	LastQuestionStartTime time.Time          `bson:"last_question_start_time,omitempty" json:"last_question_start_time,omitempty"`
	CreatedAt             time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time          `bson:"updated_at" json:"updated_at"`
	Status                PlayerStatus       `bson:"status" json:"status"`
}

func (p *Player) HasAnsweredQuestion(questionID primitive.ObjectID) bool {
	for _, answer := range p.Answers {
		if answer.QuestionID == questionID {
			return true
		}
	}

	return false
}

func (pa *PlayerAnswer) GetAnswerTime() time.Duration {
	return pa.EndTime.Sub(pa.StartTime)
}

func (ps PlayerStatus) InProgress() bool {
	return ps == PlayerStatusInProgress
}

func (ps PlayerStatus) Completed() bool {
	return ps == PlayerStatusCompleted
}

func (ps PlayerStatus) String() string {
	return []string{"Rejected", "In Progress", "Completed"}[ps]
}
