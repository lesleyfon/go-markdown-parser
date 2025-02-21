package utils

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sajari/fuzzy"
	"github.com/yuin/goldmark"
)

// findMisspelledWords finds misspelled words in a text using fuzzy matching
//
// Parameters:
//   - text: The text to check for misspelled words
//   - dictionary: A map of words to check against
//   - model: A fuzzy matching model
//
// Returns:
//   - map[string][]string: A map of misspelled words and their suggestions
//   - error: Any error encountered during processing

func findMisspelledWords(tokens []string, dictionary map[string]bool, model *fuzzy.Model) (map[string][]string, error) {

	// Make a map of misspelled words.
	misspelledWords := make(map[string][]string)

	// Check each word in the parsed text
	for _, wordToCheck := range tokens {
		wordLower := strings.ToLower(wordToCheck)

		// Only check words that don't exist in the dictionary
		_, ok := dictionary[wordLower]
		if !ok {
			// Get suggestions
			suggestions := model.Suggestions(wordToCheck, false)

			// Filter suggestions to only include close matches
			var filteredSuggestions []string
			for _, suggestion := range suggestions {
				// Calculate Levenshtein distance
				if LevenshteinDistance(wordLower, strings.ToLower(suggestion)) <= 2 {
					filteredSuggestions = append(filteredSuggestions, suggestion)
				}
			}

			// If there are suggestions, add them to the misspelled words map
			if len(filteredSuggestions) > 0 {
				misspelledWords[wordToCheck] = filteredSuggestions
			}
		}
	}

	return misspelledWords, nil
}

// convertToHTML converts markdown to HTML
//
// Parameters:
//   - contents: The markdown content as a byte slice
//
// Returns:
//   - string: The converted HTML content
func convertToHTML(contents []byte) (string, error) {

	var buffer bytes.Buffer
	err := goldmark.Convert(contents, &buffer)

	if err != nil {
		log.Printf("Markdown conversion failed: %v", err.Error())
		return "", err
	}
	return buffer.String(), nil
}

// ProcessMarkdownWithSpellCheck converts markdown content to HTML and highlights misspelled words.
// It returns the processed HTML content with spell-check markup and any error encountered.
//
// The function:
// 1. Converts markdown to HTML
// 2. Identifies misspelled words using fuzzy matching
// 3. Adds visual indicators for misspelled words with suggested corrections
//
// Parameters:
//   - contents: The markdown content as a byte slice
//   - dictionary: A map of words to check against
//   - fuzzyModel: A fuzzy matching model
//
// Returns:
//   - string: The processed HTML with spell-check markup
//   - error: Any error encountered during processing
func ProcessMarkdownWithSpellCheck(contents []byte, dictionaryMap map[string]bool, fuzzyModel *fuzzy.Model) (string, error) {
	start := time.Now()
	// logs duration of function.
	// defer func registers a function to run when the parent function returns
	defer func() {
		duration := time.Since(start)
		log.Printf("Processed markdown in %v", duration)
	}()

	// Get html contents
	htmlContents, err := convertToHTML(contents)

	if err != nil {
		return "", fmt.Errorf("markdown conversion failed: %w", err)
	}

	// Strip HTML tags and convert to plain text
	plainText := StripHTML(htmlContents)

	// // Tokenize text
	tokenizer := NewTokenizer()

	// // Tokenize text
	tokens := tokenizer.Tokenize(plainText)

	// misspelledWords, _ := findMisspelledWords(tokens, dictionaryMap, fuzzyModel)
	// Make a map of misspelled words
	misspelledWords := findMisspelledWordsParallel(tokens, dictionaryMap, fuzzyModel)

	if err != nil {
		return "", fmt.Errorf("spell check failed: %w", err)
	}

	// Process HTML and wrap misspelled words
	modifiedHTML, err := ProcessHTML(htmlContents, misspelledWords)

	if err != nil {
		// LOG Error
		log.Printf("HTML processing failed: %v", err.Error())
		return "", err
	}

	// LOG Success
	log.Printf("Spell check completed successfully.  %d misspelled words found.", len(misspelledWords))
	return modifiedHTML, nil
}

// chunkSlice splits a slice into smaller chunks of the specified size
func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// findMisspelledWordsParallel finds misspelled words in a text using fuzzy matching in parallel
//
// Parameters:
//   - words: A slice of words to check for misspelled words
//   - dictionary: A map of words to check against
//   - model: A fuzzy matching model
//
// Returns:
//   - map[string][]string: A map of misspelled words and their suggestions
func findMisspelledWordsParallel(words []string, dictionary map[string]bool, model *fuzzy.Model) map[string][]string {
	misspelledWords := make(map[string][]string)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Process words in chunks (e.g., 1000 words per goroutine)
	chunks := chunkSlice(words, 1000)

	for _, chunk := range chunks {
		wg.Add(1)
		// Use a closure to capture the current chunk
		go func(words []string) {
			// Defer the done call to ensure it runs when the function returns
			defer wg.Done()
			localMisspelled := make(map[string][]string)

			for _, word := range words {
				wordLower := strings.ToLower(word)
				if !dictionary[wordLower] {
					suggestions := model.Suggestions(word, false)
					if len(suggestions) > 0 {
						localMisspelled[word] = suggestions
					}
				}
			}

			// Lock once per chunk instead of per word
			mutex.Lock()
			for word, suggestions := range localMisspelled {
				misspelledWords[word] = suggestions
			}
			mutex.Unlock()
		}(chunk)
	}

	wg.Wait()
	return misspelledWords
}
