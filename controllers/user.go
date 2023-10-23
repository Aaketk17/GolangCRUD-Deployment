package controllers

import (
	"fmt"

	"github.com/Aaketk17/GolangCRUD-Deployment/models"
	"github.com/gofiber/fiber/v2"
)

// func UserLogin(c *fiber.Ctx) error {}

func UserSignUp(c *fiber.Ctx) error {
	var user models.User

	err := c.BodyParser(&user)
	if err != nil {
		c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Error in processing the data",
			"error":   err,
		})
		return err
	}

	fmt.Println("fefef",user)

	return err
}

// func GetUser(c *fiber.Ctx) error {}

// func GetUsers(c *fiber.Ctx) error {}

// func UpdateUsers(c *fiber.Ctx) error {}

// func DeleteUser(c *fiber.Ctx) error {}

// func checkNillErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
