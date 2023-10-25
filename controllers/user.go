package controllers

import (
	"context"
	"fmt"
	"time"

	conn "github.com/Aaketk17/GolangCRUD-Deployment/database"
	helpers "github.com/Aaketk17/GolangCRUD-Deployment/helpers"
	"github.com/Aaketk17/GolangCRUD-Deployment/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UserLogin(c *fiber.Ctx) error {
	var user models.User

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Error in processing the data",
			"error":   err,
		})
	}

	userExist, userExistErr := isUserExist(*user.Email)

	if userExistErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in finding user",
			"error":   userExistErr,
		})
	}

	if userExistErr == nil && !userExist {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"message": "User not found",
			"error":   userExistErr,
		})
	}

	userDetails, exist, userErr := findUserDetails(*user.Email)
	if userDetails == nil && !exist && userErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in finding user",
			"error":   userErr,
		})
	}

	pwdValid, deError := helpers.VerifyPassword(*userDetails.Password, *user.Password)

	if !pwdValid && deError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in decrypting password",
			"error":   deError,
		})
	}

	accessToken, tokenErr := helpers.GenerateToken(*userDetails.Email, *userDetails.Name, *userDetails.UserID)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in creating access token",
			"error":   tokenErr,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"accessToken": accessToken,
		"email":       *userDetails.Email,
		"name":        *userDetails.Name,
	})
}

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

	validate := validator.New(validator.WithRequiredStructEnabled())
	validationErr := validate.Struct(user)
	if validationErr != nil {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"message": "Error in payload",
			"error":   validationErr.Error(),
		})
	}

	if user.UserType == nil {
		value := "USER"
		user.UserType = &value
	}

	if *user.UserType == "ADMIN" {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Use users/addadmin route to add Admin users",
		})
	}

	fmt.Println(*user.UserType)

	hashedPwd, hashErr := helpers.HashPassword(*user.Password)
	if hashErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in hasing the password",
			"error":   hashErr,
		})
	}

	userExist, userExistErr := isUserExist(*user.Email)

	if userExistErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in validating user",
			"error":   userExistErr,
		})
	}

	if userExistErr == nil && userExist {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"message": "User already exist",
			"error":   userExistErr,
		})
	}

	_, insertError := conn.Connection.Exec(ctx, `insert into users (name, email, password, phone, user_type) values($1, $2, $3, $4, $5)`, *user.Name, *user.Email, hashedPwd, *user.Phone, *user.UserType)
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

func isUserExist(userEmail string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	rows, selectErr := conn.Connection.Query(ctx, "select user_id from users where email=$1", userEmail)
	if selectErr != nil {
		return true, selectErr
	}

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID)
		if err != nil {
			return true, err
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func findUserDetails(userEmail string) (*models.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	rows, selectErr := conn.Connection.Query(ctx, "select name, email, password, phone, user_id, book_id from users where email=$1", userEmail)
	if selectErr != nil {
		return nil, false, selectErr
	}

	var user models.User

	for rows.Next() {
		err := rows.Scan(&user.Name, &user.Email, &user.Password, &user.Phone, &user.UserID, &user.BookID)
		if err != nil {
			return nil, false, err
		}
	}
	return &user, true, nil
}

// func AddAdminUser() error {}

// func GetUser(c *fiber.Ctx) error {}

// func GetUsers(c *fiber.Ctx) error {}

// func UpdateUsers(c *fiber.Ctx) error {}

// func DeleteUser(c *fiber.Ctx) error {}

// func Logout(c *fiber.Ctx) error {}
