package userservice

import (
	"errors"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/validation"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Register(u *entity.User) error
	GetAll() ([]entity.User, error)
	FindByEmail(email string) (*entity.User, error)
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
	if !validation.Name(uf.Name) {
		return nil, errors.New("Name is required!")
	}

	if !validation.Password(uf.Password) {
		return nil, errors.New("Password must be greater than 6 characters.")
	}

	if !validation.Email(uf.Email) {
		return nil, errors.New("Emai is not valid.")
	}

	// TODO: uniqueness
	_, err := s.repo.FindByEmail(uf.Email)
	if err == nil {
		// does exists
		return nil, errors.New("User with this email exists")
	}

	// create user entity
	user := &entity.User{
		Name:     uf.Name,
		Email:    uf.Email,
		Mobile:   uf.Mobile,
		Password: uf.Password,
	}

	// repo store
	err = s.repo.Register(user)
	if err != nil {
		log.Println("Repo error: ", err)
	}

	return user, nil
}

func (s Service) List() ([]entity.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s Service) Login() {
	// TODO
}

func (s Service) Update() {
	// TODO
}
