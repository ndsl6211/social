package jwt

import (
	"log"
	"os"
	"strconv"
)

type authConfig struct {
	tokenSecret string
	expiredDays int
}

func NewAuthConfig() *authConfig {
	tokenSecret := os.Getenv("TOKEN_SECRET")
	expiredDays, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRED_DAYS"))
	if err != nil {
		log.Fatalln("getting environment TOKEN_EXPIRED_DAYS failed: " + err.Error())
	}

	auth := &authConfig{
		tokenSecret: tokenSecret,
		expiredDays: expiredDays,
	}

	return auth
}
