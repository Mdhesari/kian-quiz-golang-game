package main

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/service/userservice"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	url := "mongodb://michael:secret@localhost:27017/"
	clientOptions := options.Client().ApplyURI(url)

	cli, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	repo, err := mongorepo.New(cli.Database("freedom").Collection("users"))

	uf := userservice.UserForm{
		Name:     "Mahsa",
		Email:    "mahsfa@amini.ir",
		Mobile:   "+989122222222",
		Password: "123@123@123",
	}
	usersrv := userservice.New(repo)

	user, err := usersrv.Register(uf)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

	users, err := usersrv.List()
	if err != nil {
		panic(err)
	}

	fmt.Println(users)
}
