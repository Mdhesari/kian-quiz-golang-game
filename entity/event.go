package entity

type Event string

const (
	PlayersMatchedEvent Event = "matching.players_matched"
	GameStartedEvent    Event = "game.started"
	GamePlayerAnsweredEvent  Event = "game.player_answered"
)
