package gosim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVectorize(t *testing.T) {
	vocab := &Vocabulary{
		word2id:    map[string]int{"a": 1, "b": 2, "c": 3},
		nextTermId: 4,
	}

	// Case 1: updateVocab=false
	vec := vocab.Vectorize([]string{"c", "a", "a", "Z", "Z", "Z"}, false)
	assert.Equal(t, []Term{{Id: 1, Value: 2.0}, {Id: 3, Value: 1.0}}, vec)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, vocab.word2id)

	// Case 2: updateVocab=true
	vec = vocab.Vectorize([]string{"c", "a", "a", "Z", "Z", "Z"}, true)
	assert.Equal(t, []Term{{Id: 1, Value: 2.0}, {Id: 3, Value: 1.0}, {Id: 4, Value: 3}}, vec)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3, "Z": 4}, vocab.word2id)
}

// func BenchmarkVectorize(b *testing.B) {
// 	totalTokens := 0
// 	for n := 0; n < b.N; n++ {
// 		tokens := tokenize(s)
// 		totalTokens += len(tokens)
// 	}
// 	fmt.Printf("BenchmarkTokenize: totalTokens=%v\n", totalTokens)
// }
