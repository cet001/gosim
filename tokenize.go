package gosim

import (
	"strings"
	"unicode"
)

// Function definition for transforming unstructured document text into a list
// of tokens.
type Tokenize func(text string) []string

func MakeDefaultTokenizer() Tokenize {
	return func(text string) []string {
		// Pass 1: Split the string into "coarse" tokens
		tokens := strings.FieldsFunc(text, func(c rune) bool {
			return !(unicode.IsLetter(c) || unicode.IsNumber(c) || c == '\'' || c == '-')
		})

		// Pass 2: case-fold and trim non-alphanumeric characters.
		filteredTokens := make([]string, 0, len(tokens))
		for _, token := range tokens {
			token = strings.ToLower(token) // case folding
			token = strings.TrimFunc(token, func(c rune) bool {
				return !(unicode.IsLetter(c) || unicode.IsNumber(c))
			})

			// Discard single-character tokens while we're at it.
			if len(token) >= 2 {
				filteredTokens = append(filteredTokens, token)
			}
		}

		return filteredTokens
	}
}
