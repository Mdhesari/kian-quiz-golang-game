package mongorbac

import (
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"

	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	cli                  *mongorepo.MongoDB
	roleCollection       *mongo.Collection
	permissionCollection *mongo.Collection
	accessCollection     *mongo.Collection
}

func New(cli *mongorepo.MongoDB) *DB {
	return &DB{
		cli:                  cli,
		permissionCollection: cli.Conn().Collection("permissions"),
		roleCollection:       cli.Conn().Collection("roles"),
		accessCollection:     cli.Conn().Collection("access"),
	}
}
