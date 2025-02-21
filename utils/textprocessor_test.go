package utils

import (
	"testing"

	"github.com/sajari/fuzzy"
)

func BenchmarkFindMisspelledWords(b *testing.B) {
	// Setup test data
	dictionary := map[string]bool{
		"correct":  true,
		"spelling": true,
		"test":     true,
	}

	// Initialize fuzzy model
	model := fuzzy.NewModel()
	model.Train([]string{"correct", "spelling", "test"})

	// Create a large slice of words to test with
	testWords := make([]string, 1_000_000)

	for i := range testWords {
		if i%2 == 0 {
			testWords[i] = "correct"
		} else {
			testWords[i] = "incorect" // intentionally misspelled
		}
	}

	// Benchmark sequential version
	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			findMisspelledWords(testWords, dictionary, model)
		}
	})

	// Benchmark parallel version
	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			findMisspelledWordsParallel(testWords, dictionary, model)
		}
	})
}

// Regular test to ensure both functions return the same results
func TestFindMisspelledWordsEquivalence(t *testing.T) {
	// Setup test data
	dictionary := map[string]bool{
		"correct":  true,
		"spelling": true,
		"test":     true,
	}

	model := fuzzy.NewModel()
	model.Train([]string{"correct", "spelling", "test"})

	testWords := []string{
		"correct",
		"incorect", // misspelled
		"speling",  // misspelled
		"test",
	}

	// Get results from both functions
	sequential, err := findMisspelledWords(testWords, dictionary, model)
	if err != nil {
		t.Errorf("Error in sequential function: %v", err)
	}
	parallel := findMisspelledWordsParallel(testWords, dictionary, model)

	// Compare results
	if len(sequential) != len(parallel) {
		t.Errorf("Different number of misspelled words found: sequential=%d, parallel=%d",
			len(sequential), len(parallel))
	}

	// Check that all words and their suggestions match
	for word, seqSuggestions := range sequential {
		parSuggestions, exists := parallel[word]
		if !exists {
			t.Errorf("Word '%s' found in sequential but not in parallel", word)
			continue
		}

		if len(seqSuggestions) != len(parSuggestions) {
			t.Errorf("Different number of suggestions for word '%s': sequential=%d, parallel=%d",
				word, len(seqSuggestions), len(parSuggestions))
		}

		// Compare suggestions
		for i, suggestion := range seqSuggestions {
			if suggestion != parSuggestions[i] {
				t.Errorf("Different suggestion for word '%s': sequential='%s', parallel='%s'",
					word, suggestion, parSuggestions[i])
			}
		}
	}
}
