package config

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"time"
)

var defaultConfig = map[string]interface{}{
	"redis": map[string]interface{}{
		"host":     "127.0.0.1",
		"port":     6379,
		"password": "",
		"database": 0,
		"username": 0,
	},
	"application": map[string]interface{}{
		"graceful_shutdown_timeout": 5 * time.Second,
	},
	"auth": map[string]interface{}{
		"expires_in_minutes": 60,
	},
	"game": map[string]interface{}{
		"game_timeout":           15 * time.Minute,
		"max_question_timeout":   30 * time.Second,
		"max_score_per_question": entity.Score(5),
	},
}
