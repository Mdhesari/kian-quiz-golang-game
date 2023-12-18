package entity

type Player struct {
	ID        uint
	UserID    uint
	GameID    uint
	AnswerIDs []uint
	Score     int
}

type PlayerAnswer struct {
	ID       uint
	PlayerID uint
	AnswerID uint
}
