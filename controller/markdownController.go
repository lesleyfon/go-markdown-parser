package controller

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sajari/fuzzy"
	"github.com/yuin/goldmark"
	"main.go/utils"
)

// Add configuration struct
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
		FuzzyModelDepth:      4,
		FuzzyModelThreshold:  1,
	}

	// Load dictionary once
	dictionary = utils.StaticDict()
	if dictionary == nil {
		panic("Failed to load dictionary")
	}

	// Train model once
	fuzzyModel = fuzzy.NewModel()
	fuzzyModel.SetThreshold(config.FuzzyModelThreshold)
	fuzzyModel.SetDepth(config.FuzzyModelDepth)
	fuzzyModel.Train(dictionary)

	// Populate the dictionary map. This is done to avoid loading the dictionary multiple times
	for _, word := range dictionary {
		dictionaryMap[strings.ToLower(word)] = true
	}
}

func SpellCheckMarkdown() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		var buff bytes.Buffer

		// Convert markdown to html
		if err := goldmark.Convert(contents, &buff); err != nil {
			log.Printf("Markdown conversion failed: %v", err.Error())
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Markdown conversion failed: " + err.Error(),
				})
			return
		}

		// Get html contents
		htmlContents := buff.String()

		// Strip HTML tags and convert to plain text
		plainText := utils.StripHTML(htmlContents)

		// Tokenize text
		tokenizer := utils.NewTokenizer()

		// Tokenize text
		parsedText := tokenizer.Tokenize(plainText)

		// Make a map of misspelled words
		misspelledWords := make(map[string][]string)

		// Check each word in the parsed text
		for _, wordToCheck := range parsedText {
			wordLower := strings.ToLower(wordToCheck)

			// Only check words that don't exist in the dictionary
			_, ok := dictionaryMap[wordLower]
			if !ok {
				// Get suggestions
				suggestions := fuzzyModel.Suggestions(wordToCheck, false)

				// Filter suggestions to only include close matches
				var filteredSuggestions []string
				for _, suggestion := range suggestions {
					// Calculate Levenshtein distance
					if utils.LevenshteinDistance(wordLower, strings.ToLower(suggestion)) <= 2 {
						filteredSuggestions = append(filteredSuggestions, suggestion)
					}
				}

				// If there are suggestions, add them to the misspelled words map
				if len(filteredSuggestions) > 0 {
					misspelledWords[wordToCheck] = filteredSuggestions
				}
			}
		}

		// Process HTML and wrap misspelled words
		modifiedHTML, err := utils.ProcessHTML(htmlContents, misspelledWords)
		if err != nil {
			// LOG Error
			log.Printf("HTML processing failed: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "HTML processing failed: " + err.Error(),
			})
			return
		}

		// LOG Success
		log.Printf("Spell check completed successfully. %d words checked, %d misspelled words found.", len(parsedText), len(misspelledWords))

		// Return the modified HTML
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(modifiedHTML))

	}
}
