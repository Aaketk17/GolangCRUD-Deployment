package controllers

import (
	"context"
	"fmt"
	"time"

	conn "github.com/Aaketk17/GolangCRUD-Deployment/database"
	helpers "github.com/Aaketk17/GolangCRUD-Deployment/helpers"
	"github.com/Aaketk17/GolangCRUD-Deployment/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// func UserLogin(c *fiber.Ctx) error {}

func UserSignUp(c *fiber.Ctx) error {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Error in processing the data",
			"error":   err,
		})
	}
	fmt.Println("Comming....1")

	hashedPwd, hashErr := helpers.HashPassword(*user.Password)
	if hashErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in hasing the password",
			"error":   hashErr,
		})
	}
	fmt.Println("Comming....2")

	userExist, userExistErr := isUserExist(*user.Email)
	fmt.Println("Comming....3")

	if userExistErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in validating user",
			"error":   userExistErr,
		})
		// return userExistErr
	}
	fmt.Println("Comming....3")

	if userExistErr == nil && userExist {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"message": "User already exist",
			"error":   userExistErr,
		})
	}
	fmt.Println("Comming....4")

	_, insertError := conn.Connection.Exec(ctx, `insert into users (name, email, password, phone) values($1, $2, $3, $4)`, *user.Name, *user.Email, hashedPwd, *user.Phone)
	if insertError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in creating the user",
			"error":   insertError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": fmt.Sprintf("User %s created sucessfully", *user.Name),
	})
}

func isUserExist(email string) (bool, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	// defer cancel()
	fmt.Println("Comming....43", email)

	rows, queryErr := conn.Connection.Query(context.Background(), "select * from users where email=$1", email)
	if queryErr != nil {
		fmt.Println("1", queryErr)
		return true, queryErr
	}

	numbers, collectErr := pgx.CollectRows(rows, pgx.RowTo[int32])
	fmt.Println("1000", rows)

	if collectErr != nil {
		fmt.Println("2", collectErr)
		return true, collectErr
	}

	if len(numbers) > 0 {
		return true, nil
	} else {
		return false, nil
	}
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
