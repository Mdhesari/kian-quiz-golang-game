package main

import (
	"fmt"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"

	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
)

func main() {
	cli, err := mongorepo.New(mongorepo.Config{
		Username:        "michael",
		Password:        "secret",
		Port:            27017,
		Host:            "localhost",
		DBName:          "mongo",
		Migrations:      "migrations",
		DurationSeconds: 5,
	}, encrypt.Hash{})

	repo := mongouser.New(cli)

	uf := userservice.UserForm{
		Name:     "Mahsa",
		Email:    "mahsfa@amini.ir",
		Mobile:   "+989122222222",
		Password: "123@123@123",
	}
	// TODO: token should not be there!
	authSrv := authservice.New("test")
	usersrv := userservice.New(&authSrv, repo)

	res := usersrv.Register(uf)
	if len(res.Errors) > 0 {
		panic(res.Errors)
	}

	fmt.Println(res)

	users, err := usersrv.List()
	if err != nil {
		panic(err)
	}

	fmt.Println(users)
}
