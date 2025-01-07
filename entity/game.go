package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameStatus uint8

const (
	GameStatusAborted GameStatus = iota
	GameStatusInProgress
	GameStatusCompleted
)

type Game struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryID   primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	Status       GameStatus         `bson:"status" json:"status"`
	Questions    []Question         `bson:"questions,omitempty" json:"-"`
	Players      map[string]Player  `bson:"players,omitempty" json:"players,omitempty"`
	WinnerPlayer Player             `bson:"winner_player,omitempty" json:"winner_player,omitempty"`
	StartTime    time.Time          `bson:"start_time" json:"start_time"`
	EndTime      time.Time          `bson:"end_time" json:"end_time"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

func (g *Game) IsInProgress() bool {
	return g.Status == GameStatusInProgress
}

func (g *Game) IsCompleted() bool {
	return g.Status == GameStatusCompleted
}

func (g *Game) IsAborted() bool {
	return g.Status == GameStatusAborted
}

func (g *Game) GetQuestion(qId primitive.ObjectID) Question {
	for _, q := range g.Questions {
		if q.ID.Hex() == qId.Hex() {
			return q
		}
	}

	return Question{}
}
