package utils

import "fmt"

type SpellCheckError struct {
	Code    string
	Message string
	Err     error
}

func (e *SpellCheckError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

var (
	ErrInvalidFile    = &SpellCheckError{Code: "INVALID_FILE", Message: "Invalid file type"}
	ErrProcessingFile = &SpellCheckError{Code: "PROCESSING_ERROR", Message: "Error processing file"}
)
