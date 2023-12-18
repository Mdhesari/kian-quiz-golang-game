package entity

import "time"

type Game struct {
	ID          uint
	CategoryID  uint
	QuestionIDs []uint
	PlayerIDs   []uint
	WinnerID    uint
	StartTime   time.Time
	ExpiresTime time.Time
}
