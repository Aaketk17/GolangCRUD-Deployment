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

func VerifyPassword(hashedPwd string, oriPwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(oriPwd))
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error in decoding the password")
		return false, err
	}
	return true, nil
}
