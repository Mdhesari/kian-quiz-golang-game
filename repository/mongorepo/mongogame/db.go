package mongogame

import (
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"

	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	cli        *mongorepo.MongoDB
	collection *mongo.Collection
}

func New(cli *mongorepo.MongoDB) *DB {
	return &DB{
		cli:        cli,
		collection: cli.Conn().Collection("games"),
	}
}
