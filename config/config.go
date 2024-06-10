package config

import (
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
)

type Database struct {
	Migrations string           `koanf:"migrations"`
	Seeders    string           `koanf:"seeders"`
	MongoDB    mongorepo.Config `koanf:"mongodb"`
}

type JWT struct {
	Secret string `koanf:"secret"`
}

type Config struct {
	Database Database          `koanf:"database"`
	JWT      JWT               `koanf:"jwt"`
	Server   httpserver.Config `koanf:"server"`
}
