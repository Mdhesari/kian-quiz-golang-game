package userservice

import (
	"context"
	"errors"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/service/authservice"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
	IncrementScore(ctx context.Context, id primitive.ObjectID, score entity.Score) error
	FindManyById(ctx context.Context, ids []primitive.ObjectID) ([]entity.User, error)
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

func (s Service) FindMany(ctx context.Context, req param.UserFindRequest) (param.UserFindResponse, error) {
	op := "User Service: find many by ids."

	users, err := s.repo.FindManyById(ctx, req.UserIds)
	if err != nil {

		return param.UserFindResponse{}, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.UserFindResponse{
		Users: users,
	}, nil
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
		logger.L().Error("Repo error: ", zap.Error(err))

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

		logger.L().Error("Error on finding email: ", zap.Error(err))

		return &param.LoginResponse{}, richerror.New(op, err.Error()).WithKind(richerror.KindInvalid).WithErr(err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			logger.L().Error("Error on hashing password: ", zap.Error(err))
		}

		return &param.LoginResponse{}, richerror.New(op, "Credentials do not match.").WithKind(richerror.KindForbidden)
	}

	token, err := s.authSrv.GenerateToken(user, "token")
	if err != nil {
		logger.L().Error("Login service: Could not generate token.", zap.Error(err))

		return &param.LoginResponse{}, richerror.New(op, errmsg.ErrLoginTokenGenerationFailed).WithErr(err)
	}

	return &param.LoginResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s Service) IncScore(ctx context.Context, req param.UserIncrementScoreRequest) error {
	op := "User Service Increment Score"

	if err := s.repo.IncrementScore(ctx, req.UserId, req.Score); err != nil {
		logger.L().Error("Could not increment score.", zap.Error(err), zap.String("user_id", req.UserId.Hex()))

		return richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
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
