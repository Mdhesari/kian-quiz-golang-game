package userservice

import (
	"context"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/validation"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	jwt.RegisteredClaims
}

type Service struct {
	repo  Repository
	token string
}

type Repository interface {
	Register(ctx context.Context, u entity.User) (entity.User, error)
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

func (s Service) Register(uf UserForm) *param.RegisterResponse {
	res := param.RegisterResponse{
		User:   nil,
		Errors: []string{},
	}

	// TODO: validate form
	if !validation.Name(uf.Name) {
		res.Errors = append(res.Errors, "Name is required!")
	}

	if !validation.Password(uf.Password) {
		res.Errors = append(res.Errors, "Password must be greater than 6 characters.")
	}

	if !validation.Email(uf.Email) {
		res.Errors = append(res.Errors, "Email is not valid.")
	}

	if _, err := s.repo.FindByEmail(context.Background(), uf.Email); err == nil {
		res.Errors = append(res.Errors, "User with this email exists")
	}

	if len(res.Errors) > 0 {
		return &res
	}

	// create user entity
	user := entity.User{
		Name:     uf.Name,
		Email:    uf.Email,
		Mobile:   uf.Mobile,
		Password: uf.Password,
	}

	// repo store
	user, err := s.repo.Register(context.Background(), user)
	if err != nil {
		log.Println("Repo error: ", err)

		res.Errors = append(res.Errors, "Something went wrong!")
		return &res
	}

	res.User = &user
	return &res
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
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("Error on hashing password: ", err)
		}

		return &param.LoginResponse{
			Token:  "",
			Errors: []string{"Credentials do not match."},
		}
	}

	token, err := createToken(user)
	if err != nil {
		log.Println(err)

		return &param.LoginResponse{
			Token:  "",
			Errors: []string{"Something went wrong!."},
		}
	}

	return &param.LoginResponse{
		Token:  token,
		Errors: []string{},
	}
}

func (s Service) Update() {
	// TODO
}

func createToken(user *entity.User) (string, error) {
	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := Claims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}
