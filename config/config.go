package config

import (
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/delivery/grpcserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
	Game                    Game          `koanf:"game"`
}

type Game struct {
	QuestionsCount int `koanf:"questions_count"`
}

type Database struct {
	Migrations string           `koanf:"migrations"`
	Seeders    string           `koanf:"seeders"`
	MongoDB    mongorepo.Config `koanf:"mongodb"`
}

type JWT struct {
	Secret string `koanf:"secret"`
}

type Server struct {
	HttpServer httpserver.Config `koanf:"http_server"`
	GrpcServer grpcserver.Config `koanf:"grpc_server"`
}

type Config struct {
	Application Application            `koanf:"application"`
	Presence    presenceservice.Config `koanf:"presence"`
	Scheduler   scheduler.Config       `koanf:"scheduler"`
	Database    Database               `koanf:"database"`
	Server      Server                 `koanf:"server"`
	Redis       redisadapter.Config    `koanf:"redis"`
	Auth        authservice.Config     `koanf:"auth"`
	Matching    matchingservice.Config `koanf:"matching"`
}
