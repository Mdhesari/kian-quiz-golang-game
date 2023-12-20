package userservice

import (
	"errors"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Register(uf UserForm) (entity.User, error)
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

func (s Service) Register(uf UserForm) (*entity.User, error) {
	// TODO: validate form
	if len(uf.Name) < 3 {

		return nil, errors.New("Name is required!")
	}

	// TODO: uniqueness

	// repo store
	user, err := s.repo.Register(uf)
	if err != nil {
		log.Println("Repo error: ", err)
	}

	return &user, nil
}

func (s Service) Login() {
	// TODO
}

func (s Service) Update() {
	// TODO
}
