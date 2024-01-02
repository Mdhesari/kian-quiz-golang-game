package authservice

import (
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"
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
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {

		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("Could not assert claims")
}
