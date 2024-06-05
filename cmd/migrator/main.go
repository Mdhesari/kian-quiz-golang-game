package main

import (
	"fmt"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/repository/migrator"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"os"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
)

var (
	args []string
	cfg  config.Config
)

func init() {
	cfg = config.Load("config.yml")
	args = os.Args[1:]
}

func main() {
	cli, err := mongorepo.New(cfg.Database.MongoDB, encrypt.Hash{})
	if err != nil {

		panic("could not connect to mongodb.")
	}

	migrator, err := migrator.New(cli.Conn().Client(), &mongodb.Config{
		DatabaseName:         cfg.Database.MongoDB.DBName,
		MigrationsCollection: cfg.Database.MongoDB.Migrations,
		TransactionMode:      false,
		Locking:              mongodb.Locking{},
	})
	if err != nil {

		panic(err)
	}

	if len(args) > 0 && args[0] == "down" {
		migrator.Down()

		fmt.Println("Down migrations...")

		return
	}

	migrator.Up()

	fmt.Println("Up migrations...")
}
