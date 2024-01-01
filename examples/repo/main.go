package main

import (
	"fmt"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"time"

	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
)

func main() {
	cli, err := mongorepo.New(mongorepo.Config{
		Username: "michael",
		Password: "secret",
		Host:     "localhost",
		Port:     27017,
	}, 30*time.Second, encrypt.Hash{})

	repo := mongouser.New(cli)

	uf := userservice.UserForm{
		Name:     "Mahsa",
		Email:    "mahsfa@amini.ir",
		Mobile:   "+989122222222",
		Password: "123@123@123",
	}
	// TODO: token should not be there!
	usersrv := userservice.New(repo, "secret")

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
