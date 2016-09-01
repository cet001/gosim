package gosim

// NOTES:
// - http://blog.christianperone.com/2011/09/machine-learning-text-feature-extraction-tf-idf-part-i
//

import (
	"strings"
	"unicode"
)

// Function definition for transforming unstructured document text into a bag of
// terms.
type Tokenize func(text string) []string

func MakeWhitespaceTokenizer() Tokenize {
	// Determines how raw text is broken up into individual terms.
	var termSplitFn = func(c rune) bool {
		return unicode.IsSpace(c) ||
			c == '!' || c == '?' || c == ',' || c == ':' ||
			c == ';' || c == '"' || c == '|'
	}

	// Determines how each term is trimmed.
	var termTrimFn = func(c rune) bool {
		return !(unicode.IsLetter(c) || unicode.IsNumber(c))
	}

	stopWords := []string{
		"he", "than", "first", "our", "can", "they", "up", "who", "other",
		"but", "been", "one", "we", "new", "also", "their", "its", "not", "which",
		"all", "or", "said", "about", "more", "will", "have", "it", "was", "be",
		"has", "an", "are", "this", "as", "from", "by", "that", "at", "with", "is",
		"for", "on", "in", "a", "and", "of", "to", "the"}

	isStopWord := make(map[string]bool, len(stopWords))
	for _, word := range stopWords {
		isStopWord[word] = true
	}

	return func(text string) []string {
		terms := strings.FieldsFunc(text, termSplitFn)
		filteredTerms := make([]string, 0, len(terms))

		for _, term := range terms {
			term = strings.ToLower(term) // case folding
			term = strings.TrimFunc(term, termTrimFn)
			if len(term) > 0 && !isStopWord[term] {
				filteredTerms = append(filteredTerms, term)
			}
		}

		return filteredTerms
	}
}
