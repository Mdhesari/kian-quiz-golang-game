package authservice

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	Secret []byte
}

type Service struct {
	config Config
}

func New(c Config) Service {
	return Service{
		config: c,
	}
}

type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s Service) GenerateToken(user *entity.User, iss string) (string, error) {
	op := "Generate token"

	mySigningKey := []byte(s.config.Secret)

	// Create the Claims
	claims := Claims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    iss,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {

		return "", richerror.New(op, errmsg.ErrSignKey).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return ss, nil
}

func (s Service) VerifyToken(tokenString string) (*Claims, error) {
	op := "User verify token"

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		return s.config.Secret, nil
	})
	if err != nil {

		return nil, richerror.New(op, errmsg.ErrAuthorization).WithKind(richerror.KindUnAthorized)
	}

	if claims, ok := token.Claims.(*Claims); ok {

		return claims, nil
	}

	return nil, richerror.New(op, errmsg.ErrClaimAssertion).WithKind(richerror.KindUnexpected)
}
