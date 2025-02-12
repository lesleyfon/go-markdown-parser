package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"main.go/database"
	"main.go/models"
	"main.go/utils"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func MaskPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err.Error())
		return "", err
	}
	return string(bytes), nil
}

func SignUpController() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		defer cancel()

		err := c.BindJSON(&user)

		if err != nil {
			log.Println("Error binding JSON: ", err.Error())
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status":  http.StatusBadRequest,
					"message": "Error with validating user Data",
					"error":   err.Error(),
				},
			)
			return
		}

		//Check to see if users email exist
		regexpMatch := bson.M{
			"$regex": primitive.Regex{
				Pattern: *user.Email,
				Options: "i",
			},
		}
		emailCount, emailErr := userCollection.CountDocuments(ctx, bson.M{
			"email": regexpMatch,
		})

		if emailErr != nil {
			log.Println("Error checking for email: ", emailErr.Error())
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "error occurred while checking for this email",
				})
		}

		if emailCount > 0 {
			log.Println("Email already exists. emailCount: ", emailCount)
			c.JSON(
				http.StatusBadRequest, gin.H{
					"error": "Looks like this email already exists",
					"count": emailCount,
				})
			return
		}

		password, err := MaskPassword(*user.Password)

		if err != nil {
			log.Println("Error masking password: ", err.Error())
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status":  http.StatusBadRequest,
					"message": "Error with masking password",
					"error":   err.Error(),
				},
			)
			return
		}
		created_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.Password = &password
		user.Created_at = created_at
		user.Updated_at = created_at
		token, signedToken, err := utils.GenerateAllTokens(user.User_id, *user.Email)

		if err != nil {
			log.Println("Error generating token: ", err.Error())
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Error while generating token",
					"error":   err.Error(),
				},
			)
			return
		}

		user.Token = &token
		user.Refresh_token = &signedToken

		// Validate the user data
		validationError := validate.Struct(&user)
		//Check to see if data being passed meets the requirements
		if validationError != nil {
			log.Println("Error validating user data: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Error validating user data",
				"error":   validationError.Error(),
			})

			return
		}

		//To add a new user to the database
		newUser := models.User{
			ID:            user.ID,
			User_id:       user.User_id,
			Email:         user.Email,
			Password:      user.Password,
			Created_at:    user.Created_at,
			Updated_at:    user.Updated_at,
			Token:         user.Token,
			Refresh_token: user.Refresh_token,
		}

		_, err = userCollection.InsertOne(ctx, newUser)

		//Error messages
		if err != nil {
			log.Println("Error inserting user: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"data":    err.Error(),
			})
			return
		}

		log.Println("User created successfully!")
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "User created successfully!",
			"data": map[string]string{
				"user_id":       newUser.User_id,
				"token":         *newUser.Token,
				"refresh_token": *newUser.Refresh_token,
				"email":         *newUser.Email,
			},
		})
	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var user models.User
		var retrieveUser models.User

		defer cancel()

		err := c.BindJSON(&user)

		if err != nil {
			log.Println("Error binding JSON: ", err.Error())
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status":  http.StatusBadRequest,
					"message": "Error occurred while binding JSON",
					"error":   err.Error(),
				},
			)
			return
		}

		err = userCollection.FindOne(ctx,
			bson.M{
				"email": user.Email,
			}).Decode(&retrieveUser)

		if err != nil {
			log.Println("Error occurred while retrieving user: ", err.Error())
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Error occurred while retrieving user",
					"error":   err.Error(),
				},
			)
			return
		}

		passwordIsValid, msg := utils.ConfirmPassword(*user.Password, *retrieveUser.Password)

		if !passwordIsValid {
			log.Println("Invalid password: ", msg)
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"message": msg,
					"error":   "Invalid Password",
				},
			)
			return
		}

		if retrieveUser.Email == nil {
			log.Println("Oops account not found")
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status":  http.StatusBadRequest,
					"message": "Oops account not found",
				},
			)
			return
		}

		token, refreshToken, _ := utils.GenerateAllTokens(
			retrieveUser.User_id,
			*retrieveUser.Email,
		)

		updatedUser, err := utils.UpdateTokens(token, refreshToken, user.User_id)

		updatedUser.Email = retrieveUser.Email
		updatedUser.User_id = retrieveUser.User_id

		defer cancel()

		if err != nil {
			log.Println("Error updating tokens: ", err.Error())
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Error occurred while updating tokens",
					"error":   err.Error(),
				},
			)
			return
		}

		log.Println("Tokens updated successfully! login successful")
		c.JSON(
			http.StatusOK,
			gin.H{
				"email":          retrieveUser.Email,
				"userId":         retrieveUser.User_id,
				"token":          retrieveUser.Token,
				"refreshedToken": retrieveUser.Refresh_token,
			},
		)
	}
}
