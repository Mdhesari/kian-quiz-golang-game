package mongoplayer

import (
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	cli        *mongorepo.MongoDB
	collection *mongo.Collection
}

func New(cli *mongorepo.MongoDB) *MongoRepository {
	return &MongoRepository{
		cli:        cli,
		collection: cli.Conn().Collection("players"),
	}
}
