package utils

import (
	"bytes"
	"log"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Add caching for processed words
// var processedWordsCache = make(map[string]bool)

// Use sync.Map for thread-safety if needed
// var processedWordsCache sync.Map

// WrapMisspelledWordsInNode recursively processes the HTML node tree.
// For each text node, it checks if any word is misspelled and replaces it by
// wrapping that word in a <span class="misspelled-word">…</span>.
func WrapMisspelledWordsInNode(n *html.Node, misspelled map[string][]string) {
	// Regex to capture whole words if they are not a punctuation mark.
	var wordRegex = regexp.MustCompile(`\b(\w+)\b`)
	// If this is a text node, process its data.
	if n.Type == html.TextNode {
		text := n.Data
		// Skip processing if the text node is empty or only whitespace.
		if strings.TrimSpace(text) == "" {
			return
		}

		// Replace any word that is a key in the misspelled map with a <span> tag.
		newText := wordRegex.ReplaceAllStringFunc(text, func(word string) string {
			if _, exists := misspelled[word]; exists {
				return `<span class='misspelled-word bg-red-300' data-misspelled-word='` + word + `'>` + word + `</span>`
			}
			return word
		})

		// If no changes were made, there's nothing to do.
		if newText == text {
			return
		}

		// Parse the replacement HTML fragment. The parent of the current node is
		// used as the context.
		fragment, err := html.ParseFragment(strings.NewReader(newText), n.Parent)
		if err != nil {
			log.Printf("Error parsing fragment: %v", err.Error())
			return
		}

		// Insert the new nodes before the original text node.
		for _, newNode := range fragment {
			n.Parent.InsertBefore(newNode, n)
		}
		// Remove the original text node.
		n.Parent.RemoveChild(n)
	} else {
		// For non-text nodes, recursively process their children.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			WrapMisspelledWordsInNode(c, misspelled)
		}
	}
}

// ProcessHTML takes an HTML string and the misspelled words map,
// processes the document, and returns the modified HTML as a string.
func ProcessHTML(htmlStr string, misspelled map[string][]string) (string, error) {

	// Replace \n with <br>
	var charRegex = regexp.MustCompile(`\n`)
	htmlStr = charRegex.ReplaceAllString(htmlStr, "<br>")

	// Add style to html head and style the misspelled word span
	htmlStr = `
	<head>
	<script src="https://unpkg.com/@tailwindcss/browser@4"></script>
	<style>
		span.misspelled-word {
		text-decoration: underline; 
		text-decoration-style: wavy; 
		text-decoration-color: red;
		text-decoration-thickness: 1px;text-underline-offset: 3px;
		}</style>
	</head>` + htmlStr
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Printf("Error parsing html: %v", err.Error())
		return "", err
	}

	// Process the tree to wrap misspelled words.
	WrapMisspelledWordsInNode(doc, misspelled)

	// Render the modified node tree back to an HTML string.
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		log.Printf("Error rendering html: %v", err.Error())
		return "", err
	}
	return buf.String(), nil
}
