package utils

import (
	"regexp"
	"strings"
)

var (
	// Supports basic contractions and Unicode letters.
	// Precompiled regexes for performance
	tokenRegex = regexp.MustCompile(`(?i)([\p{L}]+(?:'[\p{L}]+)?|[.,!?;])`)
	alphaRegex = regexp.MustCompile(`^[\p{L}]+(?:'[\p{L}]+)?$`)
	punctRegex = regexp.MustCompile(`^[.,!?;]$`)
)

// Tokenizer struct holds the compiled regex.
type Tokenizer struct {
	re *regexp.Regexp
}

// NewTokenizer creates a new Tokenizer instance.
func NewTokenizer() *Tokenizer {
	return &Tokenizer{re: tokenRegex}
}

// Tokenize takes a string and returns a slice of cleaned tokens.
func (t *Tokenizer) Tokenize(input string) []string {
	var tokens []string
	matches := t.re.FindAllString(input, -1)

	for _, match := range matches {
		// We can optionally trim punctuation for words, if needed.
		if isAlphabetical(match) {
			// In case the token has attached punctuation, trim it.
			match = strings.TrimRight(match, ".,!?;")
		}

		// If the token is not a punctuation mark, add it to the tokens list.
		if !isPunctuation(match) {
			tokens = append(tokens, match)
		}
	}
	return tokens
}

// isAlphabetical checks if a token is a word.
func isAlphabetical(token string) bool {
	return alphaRegex.MatchString(token)
}

// isPunctuation checks if a token is a punctuation mark.
func isPunctuation(token string) bool {
	return punctRegex.MatchString(token)
}
