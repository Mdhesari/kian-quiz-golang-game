package mongorepo

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/service/userservice"
)

type Config struct {
	url string
}

type Repository struct {
	config Config
}

func New(url string) Repository {
	return Repository{
		config: Config{
			url: url,
		},
	}
}

func (r Repository) Register(uf userservice.UserForm) (entity.User, error) {
	var user entity.User

	// mongo commands

	return user, nil
}
