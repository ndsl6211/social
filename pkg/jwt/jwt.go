package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTClient interface {
	CreateToken(userId uuid.UUID) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
}

type jwtClient struct {
	tokenSecret string
	expiredAt   int64
}

func NewJWTClient(config authConfig) JWTClient {
	return &jwtClient{
		tokenSecret: config.tokenSecret,
		expiredAt:   time.Now().Add(time.Hour * 24 * time.Duration(config.expiredDays)).Unix(),
	}
}

func (jc *jwtClient) CreateToken(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userId,
		"exp": jc.expiredAt,
	})

	tokenString, err := token.SignedString([]byte(jc.tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jc *jwtClient) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jc.tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
