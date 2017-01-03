package strdist

import (
	"strings"
	"unicode"
)

// Implements the American Soundex algorithm.
//
// See https://en.wikipedia.org/wiki/Soundex
func Soundex(word string) string {
	// Based on https://github.com/dotcypress/phonetics/blob/master/soundex.go

	// Trivial case
	if len(word) == 0 {
		return "0000"
	}

	result := make([]rune, 0, len(word))
	result = append(result, unicode.ToUpper(rune(word[0])))

	// Note: we add both upper and lowercase variants of consonants to avoid the
	// expensive ToUpper(word) call.
	var code, prevCode rune
	for _, ch := range word[1:] {
		switch ch {
		case 'B', 'F', 'P', 'V', 'b', 'f', 'p', 'v':
			code = '1'
		case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z', 'c', 'g', 'j', 'k', 'q', 's', 'x', 'z':
			code = '2'
		case 'D', 'T', 'd', 't':
			code = '3'
		case 'L', 'l':
			code = '4'
		case 'M', 'N', 'm', 'n':
			code = '5'
		case 'R', 'r':
			code = '6'
		}

		if prevCode != code {
			prevCode = code
			result = append(result, code)
			if len(result) == 4 {
				break
			}
		}
	}

	if len(result) >= 4 {
		return string(result[:4])
	} else {
		return string(result) + strings.Repeat("0", 4-len(result))
	}
}

// Implments the 'Refined Soundex' algorithm (a variation on the original
// 'American Soundex' algorithm, which has fewer collisions and is typically
// more suited for spellchecking situations).
//
// See http://ntz-develop.blogspot.com/2011/03/phonetic-algorithms.html
func RefinedSoundex(word string) string {
	// Trivial case
	if len(word) == 0 {
		return "0000"
	}

	result := make([]rune, 0, len(word))
	result = append(result, unicode.ToUpper(rune(word[0])))

	// Note: we add both upper and lowercase variants of consonants to avoid the
	// expensive ToUpper(word) call.
	var code, prevCode rune
	for _, ch := range word {
		switch ch {
		case 'B', 'P', 'b', 'p':
			code = '1'
		case 'F', 'V', 'f', 'v':
			code = '2'
		case 'C', 'K', 'S', 'c', 'k', 's':
			code = '3'
		case 'G', 'J', 'g', 'j':
			code = '4'
		case 'Q', 'X', 'Z', 'q', 'x', 'z':
			code = '5'
		case 'D', 'T', 'd', 't':
			code = '6'
		case 'L', 'l':
			code = '7'
		case 'M', 'N', 'm', 'n':
			code = '8'
		case 'R', 'r':
			code = '9'
		default:
			code = '0'
		}

		if prevCode != code {
			prevCode = code
			result = append(result, code)
		}
	}

	return string(result)
}
