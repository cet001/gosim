package gosim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeWhitespaceTokenizer(t *testing.T) {
	tokenize := MakeWhitespaceTokenizer()

	assert.Equal(t, []string{}, tokenize(""))
	assert.Equal(t, []string{}, tokenize(" \n\t"))
	assert.Equal(t, []string{"foo", "bar", "baz", "foo-bar"}, tokenize(" Foo BAR \t baz!?  foo-bar\n"))

	// verify stopwords are filtered out
	assert.Equal(t, []string{"champions"}, tokenize("We are the champions."))
}

func BenchmarkTokenize(b *testing.B) {
	tokenize := MakeWhitespaceTokenizer()

	s := `
		NEW YORK—In a year that saw the release of such best-selling products as
		the Motorola RAZR 2 V8 and the wildly popular Casio XD-SW4800 handheld
		dictionary, no personal electronics product launch was more highly
		anticipated than the November 13 debut of the second-generation Microsoft
		Zune mp3 player.  The sleek new Zune, whose record-breaking sales have
		made the Zune name synonymous with "mp3 player," was so sought-after that
		thousands formed long lines outside hip, minimalist Microsoft Stores across
		the country days before the device went on sale. In Midtown Manhattan,
		the hysteria reached such a fever pitch that some were willing to pay as
		much as $200 for a spot in line.
	`

	totalTokens := 0
	for n := 0; n < b.N; n++ {
		tokens := tokenize(s)
		totalTokens += len(tokens)
	}
	fmt.Printf("BenchmarkTokenize: totalTokens=%v\n", totalTokens)
}
