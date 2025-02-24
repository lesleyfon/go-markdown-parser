package controller

import (
	"context"
	"encoding/base64"
	"go-markdown-parser/database"
	"go-markdown-parser/models"
	"go-markdown-parser/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sajari/fuzzy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var fileCollection *mongo.Collection = database.OpenCollection(database.Client, "file")

// SpellCheckConfig contains configuration parameters for spell checking
type SpellCheckConfig struct {
	LevenshteinThreshold int
	FuzzyModelDepth      int
	FuzzyModelThreshold  int
}

// Initialize a global variable for the dictionary and model
// This is done to avoid loading the dictionary and model multiple times
var (
	fuzzyModel    *fuzzy.Model
	dictionary    []string
	dictionaryMap = make(map[string]bool)
)

// init is called when the package is first loaded. Slow cold start but faster requests
// This is called only once when the server starts
func init() {
	config := SpellCheckConfig{
		LevenshteinThreshold: 2,
		FuzzyModelDepth:      1,
		FuzzyModelThreshold:  1,
	}

	// Load dictionary from file
	dictionary = utils.ImportEnglishDictionary()
	if dictionary == nil {
		panic("Failed to load dictionary")
	}

	// Configure and train model
	fuzzyModel = fuzzy.NewModel()
	fuzzyModel.SetThreshold(config.FuzzyModelThreshold)
	fuzzyModel.SetDepth(config.FuzzyModelDepth)
	fuzzyModel.Train(dictionary)

	// Populate map
	for _, word := range dictionary {
		dictionaryMap[strings.ToLower(word)] = true
	}
}

func SpellCheckMarkdown() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		file, err := c.FormFile("markdownfile")
		if err != nil {
			log.Printf("Error getting form file: %v", err.Error())
			c.JSON(http.StatusBadRequest,
				gin.H{
					"message": "Bad Request",
					"error":   err.Error(),
				})
			return
		}

		filetype := strings.Split(file.Header.Get("Content-Type"), "/")[1]

		if filetype != "markdown" {
			log.Printf("Invalid file type: %s", filetype)
			c.JSON(http.StatusBadRequest,
				gin.H{
					"message": "invalid file type. API supports only markdown files `.md`",
				})
			return
		}

		// Open markdown file
		fileContents, err := file.Open()

		if err != nil {
			log.Printf("File open failed: %v", err.Error())
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "could not open file: " + err.Error(),
				})
			return
		}
		// Close file after reading
		defer fileContents.Close()

		// Read file contents
		contents, err := io.ReadAll(fileContents)

		if err != nil {
			log.Printf("File read failed: %v", err.Error())
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Could not read file: " + err.Error(),
				})
			return
		}
		filename := file.Filename

		//TODO: MOVE ALL THIS TO AFTER ALREADY MAKING THE SPELL CHECK
		authToken := c.GetHeader("Authorization")

		// If token exist, save the db
		if authToken != "" {
			bearerToken := strings.Split(authToken, " ")[1]
			claims, _ := utils.ValidateToken(bearerToken)

			// If token is valid, save the db
			if claims != nil {
				if err := SaveMarkdownFile(ctx, filename, contents, claims.Uid); err != nil {
					log.Printf("Error saving file: %v", err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Error saving file: " + err.Error(),
					})
					return
				}
			}
		}

		// Process HTML and wrap misspelled words
		modifiedHTML, err := utils.ProcessMarkdownWithSpellCheck(contents, dictionaryMap, fuzzyModel)
		if err != nil {
			// LOG Error
			log.Printf("HTML processing failed: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "HTML processing failed: " + err.Error(),
			})
			return
		}

		// Respond with the modified HTML
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(modifiedHTML))

	}
}

// GetAllFiles returns all files for a user
// Only authenticated user can see their files
func GetAllFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		authToken := c.GetHeader("Authorization")
		// If no token, return unauthorized
		if authToken == "" {
			log.Printf("Unauthorized. Please login to continue.")

			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized. Please login to continue.",
			})
			return
		}

		bearerToken := strings.Split(authToken, " ")[1]
		claims, msg := utils.ValidateToken(bearerToken)

		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid token. Please login to continue",
				"error":   msg,
			})
			return
		}
		userId := claims.Uid

		filter := bson.M{
			"user_id": userId,
		}

		cursor, err := fileCollection.Find(ctx, filter)

		if err != nil {
			log.Printf("Error fetching files: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal server error",
			})
			return
		}

		var files []bson.M

		if err = cursor.All(ctx, &files); err != nil {
			log.Printf("Error fetching files: %v", err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Files fetched successfully",
			"files":   files,
		})
	}
}

func GetFileById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		authToken := c.GetHeader("Authorization")
		// If no token, return unauthorized
		if authToken == "" {
			log.Printf("Unauthorized. Please login to continue.")

			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized. Please login to continue.",
			})
			return
		}

		bearerToken := strings.Split(authToken, " ")[1]
		claims, msg := utils.ValidateToken(bearerToken)

		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid token. Please login to continue",
				"error":   msg,
			})
			return
		}

		fileId := c.Param("file_id")
		userId := claims.Uid

		filter := bson.M{
			"user_id": userId,
			"file_id": fileId,
		}

		var file models.File

		err := fileCollection.FindOne(ctx, filter).Decode(&file)

		if err != nil {
			log.Printf("Error fetching file: %v", err.Error())
			if err == mongo.ErrNoDocuments {
				log.Printf("File not found")
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"message": "File not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal server error",
				"error":   err.Error(),
			})
			return
		}
		// Initialize the map
		responseData := make(bson.M)

		markdownFileContents := []byte(file.File_content)
		responseData["file_content"] = file.File_content
		responseData["file_name"] = file.File_name
		responseData["user_id"] = file.User_id
		responseData["created_at"] = file.Created_at
		responseData["updated_at"] = file.Updated_at

		// Process HTML and wrap misspelled words
		modifiedHTML, err := utils.ProcessMarkdownWithSpellCheck(markdownFileContents, dictionaryMap, fuzzyModel)
		if err != nil {
			// LOG Error
			log.Printf("HTML processing failed: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "HTML processing failed: " + err.Error(),
			})
			return
		}

		// convert modifiedHTML to base64
		htmlBase64 := base64.StdEncoding.EncodeToString([]byte(modifiedHTML))
		responseData["html_content"] = htmlBase64

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "File fetched successfully",
			"file":    responseData,
		})
	}
}

// SaveMarkdownFile saves or updates a markdown file in the database
func SaveMarkdownFile(ctx context.Context, filename string, contents []byte, userId string) error {
	fileFilter := bson.M{
		"file_name": filename,
		"user_id":   userId,
	}

	// Prepare the file document
	now := time.Now()
	fileDoc := bson.M{
		"file_name":    filename,
		"user_id":      userId,
		"file_content": string(contents),
		"updated_at":   now,
	}

	// Update the file
	result, err := fileCollection.UpdateOne(
		ctx,
		fileFilter,
		bson.M{"$set": fileDoc},
	)

	if err != nil {
		log.Printf("Error occurred while updating file: %v", err.Error())
		return err
	}

	// If no document was updated, create new one
	if result.MatchedCount == 0 {
		docId := primitive.NewObjectID()
		fileDoc["_id"] = docId
		fileDoc["file_id"] = docId.Hex()
		fileDoc["created_at"] = now

		_, err := fileCollection.InsertOne(ctx, fileDoc)
		// If error, return error
		if err != nil {
			log.Printf("Failed to create new file: %v", err.Error())
			return err
		}
		log.Printf("New file created successfully")
	} else {
		log.Printf("File updated successfully")
	}

	return nil
}
