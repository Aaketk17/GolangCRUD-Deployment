package helpers

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, name string, user_id int) (string, error) {
	tokenSecret := os.Getenv("JWT_SECRET")
	claims := UserClaims{
		Email:  email,
		Name:   name,
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 6 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(tokenSecret))
	if err != nil {
		log.Panic(err)
		return "", err
	}

	return token, nil
}
