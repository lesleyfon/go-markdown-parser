package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ConfirmPassword(userPassword string, passwordEntered string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordEntered), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		log.Println("Error comparing password: ", err.Error())
		msg = "Looks like you entered a wrong password"
		check = false
	}
	return check, msg
}
