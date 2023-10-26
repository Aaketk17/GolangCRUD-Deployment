package controllers

import (
	"context"
	"fmt"
	"strconv"
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

	userDetails, exist, userErr := findUserDetailsByEmail(*user.Email)
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

	accessToken, tokenErr := helpers.GenerateToken(*userDetails.Email, *userDetails.Name, *userDetails.UserID, *userDetails.UserType)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in creating access token",
			"error":   tokenErr,
		})
	}

	tokenCookie := new(fiber.Cookie)
	tokenCookie.Name = "accessToken"
	tokenCookie.Value = accessToken
	tokenCookie.HTTPOnly = true
	tokenCookie.Secure = true
	tokenCookie.Expires = time.Now().Add(24 * time.Hour)
	tokenCookie.SameSite = "Strict"

	c.Cookie(tokenCookie)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"email":   *userDetails.Email,
		"name":    *userDetails.Name,
		"message": "Login success",
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
		return false, selectErr
	}

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID)
		if err != nil {
			return false, err
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func findUserDetailsByEmail(userEmail string) (*models.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	rows, selectErr := conn.Connection.Query(ctx, "select name, email, password, phone, user_id, user_type from users where email=$1", userEmail)
	if selectErr != nil {
		return nil, false, selectErr
	}

	var user models.User

	for rows.Next() {
		err := rows.Scan(&user.Name, &user.Email, &user.Password, &user.Phone, &user.UserID, &user.UserType)
		if err != nil {
			return nil, false, err
		}
	}
	if *user.Email == "" || *user.UserID == 0 {
		return nil, false, nil
	}
	return &user, true, nil
}

func findUserDetailsByUserId(userID int) (*models.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	rows, selectErr := conn.Connection.Query(ctx, "select name, email, password, phone, user_type, user_id from users where user_id=$1", userID)
	if selectErr != nil {
		return nil, false, selectErr
	}

	var user models.User

	for rows.Next() {
		err := rows.Scan(&user.Name, &user.Email, &user.Password, &user.Phone, &user.UserType, &user.UserID)
		if err != nil {
			return nil, false, err
		}
	}

	if user.Email == nil || user.UserID == nil {
		return nil, false, nil
	}

	return &user, true, nil
}

func AddAdminUser(c *fiber.Ctx) error {
	token, tokenErr := helpers.GetTokenFromCookies(c)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in getting token from cookies",
			"error":   tokenErr,
		})
	}

	claims, errMsg, valid := helpers.GetTokenClaims(token)
	if !valid {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": errMsg,
		})
	}

	if claims.UserType != "ADMIN" {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Only admin users can add another admin user",
		})
	}

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

	if *user.UserType == "USER" {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Use /signup route to add normal user",
		})
	}

	if user.UserType == nil {
		value := "ADMIN"
		user.UserType = &value
	}

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
		"message": fmt.Sprintf("Admin user %s created sucessfully", *user.Name),
	})
}

func GetUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	token, tokenErr := helpers.GetTokenFromCookies(c)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in getting token from cookies",
			"error":   tokenErr,
		})
	}

	claims, errMsg, valid := helpers.GetTokenClaims(token)
	if !valid {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": errMsg,
		})
	}

	userIdParamsInt, convErr := strconv.Atoi(userId)
	if convErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in converting userId string to integer",
			"error":   convErr,
		})
	}
	if claims.UserType == "USER" && userIdParamsInt != claims.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Normal users can't access other users account",
		})
	}

	userDetails, userExists, userErr := findUserDetailsByUserId(userIdParamsInt)
	fmt.Println(userExists, "ererer", userErr)
	if !userExists && userErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": fmt.Sprintf("Error in finding user by userid %s", userId),
			"error":   userErr,
		})
	}

	if userDetails == nil && !userExists {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": fmt.Sprintf("User with userID %s not found", userId),
			"error":   userErr,
		})
	}
	return c.Status(fiber.StatusOK).JSON(userDetails)
}

func GetUsers(c *fiber.Ctx) error {
	token, tokenErr := helpers.GetTokenFromCookies(c)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in getting token from cookies",
			"error":   tokenErr,
		})
	}

	claims, errMsg, valid := helpers.GetTokenClaims(token)
	if !valid {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": errMsg,
		})
	}

	if claims.UserType == "USER" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Normal users can't access other users account",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	rows, selectErr := conn.Connection.Query(ctx, "select name, email, password, phone, user_type, user_id from users")
	if selectErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error in getting data from database",
			"error":   selectErr,
		})
	}

	var user models.User
	var usersData []models.User

	for rows.Next() {
		rowErr := rows.Scan(&user.Name, &user.Email, &user.Password, &user.Phone, &user.UserType, &user.UserID)
		if rowErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"message": "Error in scanning rows",
				"error":   rowErr,
			})
		}

		usersData = append(usersData, user)
	}

	if len(usersData) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "No users in the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(usersData)
}

// func UpdateUsers(c *fiber.Ctx) error {}

// func DeleteUser(c *fiber.Ctx) error {}

// func Logout(c *fiber.Ctx) error {}
