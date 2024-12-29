package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryID   primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	Questions    []Question         `bson:"questions,omitempty" json:"questions,omitempty"`
	Players      []Player           `bson:"players,omitempty" json:"players,omitempty"`
	WinnerPlayer Player             `bson:"winner_player,omitempty" json:"winner_player,omitempty"`
	StartTime    time.Time          `bson:"start_time" json:"start_time"`
	EndTime      time.Time          `bson:"end_time" json:"end_time"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
