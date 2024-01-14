package authservice

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	secret []byte
}

func New(secret string) Service {
	return Service{
		secret: []byte(secret),
	}
}

type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s Service) GenerateToken(user *entity.User, iss string) (string, error) {
	mySigningKey := []byte(s.secret)

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
		return "", err
	}

	return ss, nil
}

func (s Service) VerifyToken(tokenString string) (*Claims, error) {
	op := "User verify token"

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {

		return nil, richerror.New(op, errmsg.ErrAuthorization).WithKind(richerror.KindUnAthorized)
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, richerror.New(op, errmsg.ErrClaimAssertion).WithKind(richerror.KindUnexpected)
}
