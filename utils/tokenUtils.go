package utils

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/database"
	"main.go/models"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

type JwtSignedDetails struct {
	Email     string
	Name      string
	Username  string
	Uid       string
	User_type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var secretKeyBytes = []byte(SECRET_KEY)

// GenerateAllTokens generates a new token and refresh token for a user
func GenerateAllTokens(
	uid string,
	email string,
) (
	signedToken string,
	signedRefreshToken string,
	err error,
) {
	// Create claims once and reuse
	now := time.Now().Local()
	standardClaims := jwt.StandardClaims{
		ExpiresAt: now.Add(time.Hour * 12).Unix(),
	}

	refreshStandardClaims := jwt.StandardClaims{
		ExpiresAt: now.Add(time.Hour * 100).Unix(),
	}

	// Create the claims for the token
	claims := &JwtSignedDetails{
		Uid:            uid,
		Email:          email,
		StandardClaims: standardClaims,
	}

	// Create the claims for the refresh token
	refreshClaims := &JwtSignedDetails{
		StandardClaims: refreshStandardClaims,
	}

	// Create both tokens in parallel using goroutines
	var tokenErr, refreshTokenErr error
	var token, refreshToken string

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		token, tokenErr = jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			claims,
		).SignedString(secretKeyBytes)
	}()

	go func() {
		defer wg.Done()
		refreshToken, refreshTokenErr = jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			refreshClaims,
		).SignedString(secretKeyBytes)
	}()

	wg.Wait()

	// Check for errors
	if tokenErr != nil {
		return "", "", tokenErr
	}
	if refreshTokenErr != nil {
		return "", "", refreshTokenErr
	}

	return token, refreshToken, nil
}

func UpdateTokens(signedToken string, signedRefreshedToken string, userId string) (*models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Initialize the update document
	updateTokenDocs := bson.D{
		{
			Key:   "token",
			Value: signedToken,
		},
		{
			Key:   "refreshedToken",
			Value: signedRefreshedToken,
		},
		{
			Key:   "updated_at",
			Value: time.Now(),
		},
	}

	// Define filter based on the userId
	filter := bson.M{"user_id": userId}
	returnDocument := options.After

	// Specify options for upsert and to return the updated document
	upsert := true
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnDocument, // Return the updated document
	}

	// Perform the update operation and get the updated document back
	var updatedUser models.User
	err := userCollection.FindOneAndUpdate(
		ctx,
		filter,
		bson.D{
			{
				Key:   "$set",
				Value: updateTokenDocs,
			},
		},
		&opt,
	).Decode(&updatedUser) // Decode the result into updatedUser

	if err != nil {
		log.Printf("Error updating token for user %s: %v", userId, err.Error())
		return nil, err
	}

	return &updatedUser, nil
}
