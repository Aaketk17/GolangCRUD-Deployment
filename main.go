package main

import (
	"context"
	"fmt"
	"os"

	con "github.com/Aaketk17/GolangCRUD-Deployment/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("This is golang with postgress CRUD project")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error in loading .env file")
	}

	app := fiber.New()
	app.Get("/home", func(c *fiber.Ctx) error {
		_, err := con.Connection.Query(context.Background(), `select * from students`)
		return err

	})
	app.Listen(":" + os.Getenv("PORT"))
}
