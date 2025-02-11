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

func StaticDict() []string {

	correct := []string{
		"A", "a",
		"Above", "above",
		"Age", "age",
		"Alice", "alice",
		"And", "and",
		"Are", "are",
		"Basic", "basic",
		"Be", "be",
		"Below", "below",
		"Check", "check",
		"Checker", "checker",
		"Code", "code",
		"Com", "com",
		"Console", "console",
		"Corrections", "corrections",
		"Correctly", "correctly",
		"Custom", "custom",
		"Didn't", "didn't",
		"Doe", "doe",
		"Demonstration", "demonstration",
		"Designed", "designed",
		"Details", "details",
		"Document", "document",
		"Doing", "doing",
		"Errors", "errors",
		"Example", "example",
		"First", "first",
		"File", "file",
		"Focus", "focus",
		"For", "for",
		"Function", "function",
		"Goal", "goal",
		"Grammatical", "grammatical",
		"Greet", "greet",
		"Greets", "greets",
		"Hello", "hello",
		"How", "how",
		"I", "i",
		"If", "if",
		"Identify", "identify",
		"In", "in",
		"Incorrectly", "incorrectly",
		"Information", "information",
		"Intended", "intended",
		"Introduction", "introduction",
		"Is", "is",
		"John", "john",
		"John doe", "john doe",
		"Markdown", "markdown",
		"May", "may",
		"Message", "message",
		"My", "my",
		"Name", "name",
		"Occur", "occur",
		"Of", "of",
		"Log", "log",
		"On", "on",
		"One", "one",
		"Order", "order",
		"Provide", "provide",
		"Purpose", "purpose",
		"Proper", "proper",
		"Real-life", "real-life",
		"Real", "real",
		"Intended", "intended",
		"Life", "life",
		"Sample", "sample",
		"Scenario", "scenario",
		"Scenarios", "scenarios",
		"Section", "section",
		"Several", "several",
		"Showcase", "showcase",
		"Simple", "simple",
		"Simulate", "simulate",
		"Spell", "spell",
		"Spelling", "spelling",
		"Spelt", "spelt",
		"Suggest", "suggest",
		"Test", "test",
		"Testing", "testing",
		"The", "the",
		"This", "this",
		"Three", "three",
		"To", "to",
		"Today", "today",
		"Used", "used",
		"User", "user",
		"Various", "various",
		"Where", "where",
		"While", "while",
		"With", "with",
		"Writing", "writing",
		"You", "you",
		"Words", "words",
		"Simulate", "simulate",
		"Simulates", "simulates",
		"Simulated", "simulated",
		"Simulating", "simulating",
		"Simulations", "simulations",
		"Simulations", "simulations",
		"Few", "few",
		"Example", "example",
		"Examples", "examples",
		"Exampled", "exampled",
		"Exampling", "exampling",
		"Common", "common",
		"Commonly", "commonly",
		"Block", "block",
		"Blocks", "blocks",
		"Blocked", "blocked",
		"Blocking", "blocking",
		"It", "it",
		"Is", "is",
		"I'm", "i'm",
		"I've", "i've",
		"I'll", "i'll",
		"I'd", "i'd",
		"I'm", "i'm",
		"Email", "email",
		"Emails", "emails",
		"can", "Can",
	}
	return correct
}
