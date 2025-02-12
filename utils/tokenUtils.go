package utils

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtSignedDetails struct {
	Email     string
	Name      string
	Username  string
	Uid       string
	User_type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// GenerateAllTokens generates a new token and refresh token for a user
func GenerateAllTokens(
	uid string,
	email string,
) (
	signedToken string,
	signedRefreshToken string,
	err error,
) {

	// Create the claims for the token
	claims := &JwtSignedDetails{
		Uid:   uid,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(12)).Unix(),
		},
	}

	// Create the claims for the refresh token
	refreshClaims := &JwtSignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(100)).Unix(),
		},
	}

	// Create the token
	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString(
		[]byte(SECRET_KEY),
	)

	// Create the refresh token
	refreshToken, refreshTokenErr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		// Log the error
		log.Panic(err)

		// Return the error
		return
	}
	if refreshTokenErr != nil {
		// Log the error
		log.Panic(refreshTokenErr)

		// Return the error
		return
	}

	// Return the token and refresh token
	return token, refreshToken, err

}
