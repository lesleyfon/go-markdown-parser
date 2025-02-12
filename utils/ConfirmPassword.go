package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Default cost is 10, but we can adjust based on our needs
// Lower cost = faster but less secure
// Higher cost = slower but more secure
const BcryptCost = 10

func ConfirmPassword(userPassword string, hashedPassword string) (bool, string) {

	// Compare the passwords
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))

	if err != nil {
		log.Println("Error comparing password: ", err.Error())
		return false, "Invalid password"
	}

	return true, ""
}

func HashPassword(password string) (string, error) {
	// Generate password hash
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
