package userservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/validation"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

type Service struct {
	repo  Repository
	token string
}

type Repository interface {
	Register(ctx context.Context, u *entity.User) error
	GetAll(ctx context.Context) ([]entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type UserForm struct {
	Name     string
	Email    string
	Mobile   string
	Password string
}

func New(repo Repository, token string) Service {
	return Service{repo: repo, token: token}
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

	_, err := s.repo.FindByEmail(context.Background(), uf.Email)
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
	err = s.repo.Register(context.Background(), user)
	if err != nil {
		log.Println("Repo error: ", err)
	}

	return user, nil
}

func (s Service) List() ([]entity.User, error) {
	users, err := s.repo.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s Service) Login(req param.LoginRequest) *param.LoginResponse {
	user, err := s.repo.FindByEmail(context.Background(), req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return &param.LoginResponse{
				Token:  "",
				Errors: []string{"Credentials do not match."},
			}
		}

		log.Println("Error on finding email: ", err)

		return &param.LoginResponse{
			Token:  "",
			Errors: []string{err.Error()},
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println("Error on hashing password: ", err)

		return &param.LoginResponse{
			Token:  "",
			Errors: []string{err.Error()},
		}
	}

	token, err := jwt.ParseWithClaims(s.token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		log.Println("Error on hashing password: ", err)

		return &param.LoginResponse{
			Token:  "",
			Errors: []string{err.Error()},
		}
	} else if claims, ok := token.Claims.(*Claims); ok {
		fmt.Println(claims.Foo, claims.RegisteredClaims.Issuer)
	} else {
		msg := "unknown claims type, cannot proceed"
		log.Println(msg)

		return &param.LoginResponse{
			Token:  "",
			Errors: []string{msg},
		}
	}

	return &param.LoginResponse{
		Token:  "YOUR_TOKEN",
		Errors: []string{},
	}
}

func (s Service) Update() {
	// TODO
}
