package helpers

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	hashValue, hashErr := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	if hashErr != nil {
		log.Fatal(hashErr)
		fmt.Println("Error in hashing the password", hashErr)
		return "", hashErr
	}
	return string(hashValue), nil
}

// func DecryptPassword(hashPwd string) string {}
