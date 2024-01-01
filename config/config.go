package config

import "mdhesari/kian-quiz-golang-game/repository/mongorepo"

type Database struct {
	MongoDB mongorepo.Config `koanf:"mongodb"`
}

type JWT struct {
	Secret string `koanf:"secret"`
}

type Config struct {
	Database Database `koanf:"database"`
	JWT      JWT      `koanf:"jwt"`
}
