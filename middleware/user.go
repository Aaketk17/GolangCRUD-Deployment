package middleware

import (
	"context"
	"os"
	"time"

	conn "github.com/Aaketk17/GolangCRUD-Deployment/database"
	"github.com/Aaketk17/GolangCRUD-Deployment/models"
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

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	rows, selectErr := conn.Connection.Query(ctx, "select id, token, updated_at from invalidtokens where token=$1", t.Value)
	if selectErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in getting existing tokens from DB for validation",
			"error":   selectErr.Error(),
		})
	}

	var tokenModel models.InvalidToken

	for rows.Next() {
		err := rows.Scan(&tokenModel.Id, &tokenModel.Token, &tokenModel.UpdatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"message": "Error in scanning tokens from DB",
				"error":   err.Error(),
			})
		}
	}

	if tokenModel.Id != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "User already logged out, login again",
		})
	}
	return c.Next()
}
