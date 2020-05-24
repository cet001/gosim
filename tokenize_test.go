package gosim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleTokenize() {
	tokenize := MakeDefaultTokenizer()
	tokens := tokenize("Let's go bowling!")

	for i, token := range tokens {
		fmt.Printf("token %v: %v\n", i, token)
	}
	// Output:
	// token 0: let's
	// token 1: go
	// token 2: bowling
}

func TestMakeDefaultTokenizer(t *testing.T) {
	tokenize := MakeDefaultTokenizer()

	// Basic examples
	assert.Equal(t, []string{"mom's", "and", "dad's"}, tokenize("Mom's and Dad's"))
	assert.Equal(t, []string{"foo", "bar", "baz", "foo-bar"}, tokenize(" Foo BAR \t baz!?  foo-bar\n"))

	// These input strings should not produce any tokens.
	assert.Equal(t, []string{}, tokenize(""))
	assert.Equal(t, []string{}, tokenize(" \n\t"))

	// Verify single-character tokens are filtered out
	assert.Equal(t, []string{"aa", "aaa"}, tokenize("a aa aaa"))

	// Verify single- and double-quoted strings are "de-quoted"
	assert.Equal(t, []string{"one", "two", "three", "four"}, tokenize(`one "two" '''three''' 'four'`))
}

func BenchmarkTokenize(b *testing.B) {
	tokenize := MakeDefaultTokenizer()

	s := `
		NEW YORK—In a year that saw the release of such best-selling products as
		the Motorola RAZR 2 V8 and the wildly popular Casio XD-SW4800 handheld
		dictionary, no personal electronics product launch was more highly
		anticipated than the November 13 debut of the second-generation Microsoft
		Zune mp3 player.  The sleek new Zune, whose record-breaking sales have
		made the Zune name synonymous with "mp3 player", was so sought-after that
		thousands formed long lines outside hip, minimalist Microsoft Stores across
		the country days before the device went on sale. In Midtown Manhattan,
		the hysteria reached such a fever pitch that some were willing to pay as
		much as $200 for a spot in line.  That's it.
	`

	totalTokens := 0
	for n := 0; n < b.N; n++ {
		tokens := tokenize(s)
		totalTokens += len(tokens)
	}
	fmt.Printf("BenchmarkTokenize: totalTokens=%v\n", totalTokens)
}
