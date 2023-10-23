package main

import (
	"fmt"
	"os"

	routes "github.com/Aaketk17/GolangCRUD-Deployment/routes"
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

	// routes.BookRoutes(app)
	routes.UserRoutes(app)

	app.Listen(":" + os.Getenv("PORT"))
}
