package config

import "mdhesari/kian-quiz-golang-game/repository/mongorepo"

type Database struct {
	MongoDB mongorepo.Config `koanf:"mongodb"`
}

type Config struct {
	Database Database `koanf:"database"`
}