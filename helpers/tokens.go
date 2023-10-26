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

func ValidateToken(signedToken string) (claims *UserClaims, msg string) {
	var SECRET_KEY string = os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(
		signedToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		msg = "Token is not valid"
		msg = err.Error()
		return
	}
	return claims, msg

}
