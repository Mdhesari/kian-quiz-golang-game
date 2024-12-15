package entity

type Event string

const (
	UsersMatchedEvent Event = "matching.users_matched"
	GameStartedEvent  Event = "game.started"
)