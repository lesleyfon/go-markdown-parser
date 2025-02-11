package utils

import "regexp"

/*
*
HTML Stripping:

	Use regex to strip HTML tags from a string.
*/
func StripHTML(input string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllLiteralString(input, " ")
}
