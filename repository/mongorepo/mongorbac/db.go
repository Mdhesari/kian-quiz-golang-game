package mongorbac

import (
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
)

type DB struct {
	cli *mongorepo.MongoDB
}

func New(cli *mongorepo.MongoDB) *DB {
	return &DB{
		cli: cli,
	}
}
