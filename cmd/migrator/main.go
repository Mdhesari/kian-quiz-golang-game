package main

import (
	"flag"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/repository/migrator"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

var (
	seeds *bool
	down  *bool
	cfg   config.Config
)

func init() {
	cfg = config.Load("config.yml")
	seeds = flag.Bool("seeds", false, "Runs seeders.")
	down = flag.Bool("down", false, "Down migrations.")

	flag.Parse()
}

func main() {
	cli := mongorepo.New(cfg.Database.MongoDB)
	mongoConf := &mongodb.Config{
		DatabaseName:         cfg.Database.MongoDB.DBName,
		MigrationsCollection: cfg.Database.MongoDB.Migrations,
		TransactionMode:      false,
		Locking:              mongodb.Locking{},
	}
	mgrt, err := migrator.New(cli.Conn().Client(), mongoConf, cfg.Database.Migrations)
	if err != nil {

		panic(err)
	}

	if *down {
		logger.L().Info("down")
		err := mgrt.Down()
		if err != nil {
			panic(err)
		}

		logger.L().Info("Down migrations...")

		return
	}

	err = mgrt.Up()
	if err != nil {
		panic(err)
	}

	logger.L().Info("Migrations are run successfuly.")

	if *seeds {
		logger.L().Info("Running seeders...")

		mongoConf.MigrationsCollection = "seeders"
		seeder, err := migrator.New(cli.Conn().Client(), mongoConf, cfg.Database.Seeders)
		if err != nil {
			logger.L().Error("Seeders Error: ", zap.Error(err))
		}

		err = seeder.Up()
		if err != nil {
			logger.L().Error("Seeders UP Error: ", zap.Error(err))
		}

		logger.L().Info("Seeders are run successfuly.")
	}
}
