package userservice

import (
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Register(uf UserForm) (entity.User, error)
	Login()
	Update()
}

type UserForm struct {
	Name     string
	Email    string
	Mobile   string
	Password string
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Register(uf UserForm) entity.User {
	// TODO: validate form

	// TODO: uniqueness

	// repo store
	user, err := s.repo.Register(uf)
	if err != nil {
		log.Println("Repo error: ", err)
	}

	return user
}

func (s Service) Login() {
	// TODO
}

func (s Service) Update() {
	// TODO
}
