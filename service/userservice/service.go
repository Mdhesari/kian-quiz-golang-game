package userservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/service/authservice"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("Could not find the entity.")

type Service struct {
	authSrv *authservice.Service
	repo    Repository
}

type Repository interface {
	Register(ctx context.Context, u entity.User) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(id primitive.ObjectID) (*entity.User, error)
}

type UserForm struct {
	Name     string
	Email    string
	Mobile   string
	Password string
	RoleID   *primitive.ObjectID
}

func New(authSrv *authservice.Service, repo Repository) Service {
	return Service{authSrv: authSrv, repo: repo}
}

func (s Service) Register(uf UserForm) (*param.RegisterResponse, error) {
	op := "User Register"

	password, err := bcrypt.GenerateFromPassword([]byte(uf.Password), bcrypt.DefaultCost)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	// TODO - Avatar generator
	avatar := fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", uf.Email)
	user := entity.User{
		Name:     uf.Name,
		Email:    uf.Email,
		Mobile:   uf.Mobile,
		Password: password,
		RoleID:   uf.RoleID,
		Avatar:   avatar,
	}

	user, err = s.repo.Register(context.Background(), user)
	if err != nil {
		log.Println("Repo error: ", err)

		return nil, richerror.New(op, errmsg.ErrInternalServer).WithKind(richerror.KindUnexpected)
	}

	return &param.RegisterResponse{
		User: &user,
	}, nil
}

func (s Service) List() ([]entity.User, error) {
	users, err := s.repo.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s Service) Login(req param.LoginRequest) (*param.LoginResponse, error) {
	op := "User Service Login"

	user, err := s.repo.FindByEmail(context.Background(), req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return &param.LoginResponse{
				Token: "",
			}, richerror.New(op, "Credentials do not match.").WithErr(err).WithKind(richerror.KindForbidden)
		}

		log.Println("Error on finding email: ", err)

		return &param.LoginResponse{}, richerror.New(op, err.Error()).WithKind(richerror.KindInvalid).WithErr(err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("Error on hashing password: ", err)
		}

		return &param.LoginResponse{}, richerror.New(op, "Credentials do not match.").WithKind(richerror.KindForbidden)
	}

	token, err := s.authSrv.GenerateToken(user, "token")
	if err != nil {
		log.Println(err)

		return &param.LoginResponse{}, richerror.New(op, "Something went wrong!.").WithErr(err)
	}

	return &param.LoginResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s Service) Update() {
	// TODO
}

func (s Service) GetByID(id primitive.ObjectID) (param.ProfileResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return param.ProfileResponse{}, ErrNotFound
		}

		return param.ProfileResponse{}, err
	}

	return param.ProfileResponse{
		User: user,
	}, nil
}
