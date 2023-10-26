package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type tokenFromCookie struct {
	Value string `cookie:"accessToken"`
}

func UserAuthMiddleware(c *fiber.Ctx) error {

	t := new(tokenFromCookie)
	cookieErr := c.CookieParser(t)
	if cookieErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in getting token from cookies",
			"error":   cookieErr,
		})
	}

	var SECRET_KEY string = os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(
		t.Value,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Unauthorized access",
			"error":   err.Error(),
		})
	}
	return c.Next()
}
