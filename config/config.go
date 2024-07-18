package config

import (
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type Database struct {
	Migrations string           `koanf:"migrations"`
	Seeders    string           `koanf:"seeders"`
	MongoDB    mongorepo.Config `koanf:"mongodb"`
}

type JWT struct {
	Secret string `koanf:"secret"`
}

type Config struct {
	Application Application            `koanf:"application"`
	Presence    presenceservice.Config `koanf:"presence"`
	Scheduler   scheduler.Config       `koanf:"scheduler"`
	Database    Database               `koanf:"database"`
	JWT         JWT                    `koanf:"jwt"`
	Server      httpserver.Config      `koanf:"server"`
	Redis       redisadapter.Config    `koanf:"redis"`
}
