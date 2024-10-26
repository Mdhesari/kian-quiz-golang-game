package main

import (
	"flag"
	"fmt"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/repository/migrator"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/gommon/log"
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
	cli, err := mongorepo.New(cfg.Database.MongoDB)
	if err != nil {

		panic("could not connect to mongodb.")
	}

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
		err := mgrt.Down()
		if err != nil {
			panic(err)
		}

		fmt.Println("Down migrations...")

		return
	}

	err = mgrt.Up()
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrations are run successfuly.")

	if *seeds {
		fmt.Println("Running seeders...")

		mongoConf.MigrationsCollection = "seeders"
		seeder, err := migrator.New(cli.Conn().Client(), mongoConf, cfg.Database.Seeders)
		if err != nil {
			log.Fatal("Seeders Error: ", err)
		}

		err = seeder.Up()
		if err != nil {
			log.Fatal("Seeders UP Error: ", err)
		}

		fmt.Println("Seeders are run successfuly.")
	}
}
