package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Add support for custom dictionaries
type Dictionary struct {
	Words   map[string]bool
	Custom  map[string]bool
	Ignored map[string]bool
}

// Add methods to manage dictionary
func (d *Dictionary) AddCustomWord(word string) {
	d.Custom[strings.ToLower(word)] = true
}

// Ignore word
func (d *Dictionary) IgnoreWord(word string) {
	d.Ignored[strings.ToLower(word)] = true
}

func ImportEnglishDictionary() []string {
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting working directory: %v\n", err)
		return nil
	}

	// Construct path to dictionary file
	// This assumes the data directory is in the project root
	filePath := filepath.Join(wd, "data", "dictionary.txt")

	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return nil
	}
	defer readFile.Close()

	var dictionary []string

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		trimmedStr := strings.Join(strings.Fields(line), " ")
		splitStr := strings.Split(trimmedStr, " ")

		if len(splitStr) > 0 {
			dictionary = append(dictionary, splitStr[0])
		}
	}
	return dictionary
}
