package helpers

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

type tokenFromCookie struct {
	Value string `cookie:"accessToken"`
}

func GenerateToken(email string, name string, user_id int, user_type string) (string, error) {
	tokenSecret := os.Getenv("JWT_SECRET")
	claims := UserClaims{
		Email:    email,
		Name:     name,
		UserID:   user_id,
		UserType: user_type,
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

func GetTokenClaims(signedToken string) (claims *UserClaims, msg string, valid bool) {
	var SECRET_KEY string = os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(
		signedToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		valid = false
		return
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		msg = "Token is not valid :" + err.Error()
		valid = false
		return
	}
	valid = true
	return claims, msg, valid
}

func GetTokenFromCookies(c *fiber.Ctx) (token string, err error) {
	t := new(tokenFromCookie)
	cookieErr := c.CookieParser(t)
	if cookieErr != nil {
		return "", cookieErr
	}

	return t.Value, nil
}
