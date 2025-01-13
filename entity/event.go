package entity

type Event string

const (
	PlayersMatchedEvent     Event = "matching.players_matched"
	GameStartedEvent        Event = "game.started"
	GameStatusFinishedEvent Event = "game.finished"
	PlayerAnsweredEvent     Event = "player.player_answered"
	PlayerFinishedEvent     Event = "player.finished"
)
