package gosim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVocabulary_BasicUsage(t *testing.T) {
	vocab := NewVocabulary()
	assert.Equal(t, 0, vocab.Size())

	termvec := vocab.Vectorize([]string{"b", "a", "c", "b", "a", "a"}, true)

	// Verify terms are ordered by Id
	for i := 0; i < len(termvec)-1; i++ {
		assert.True(t, termvec[i].Id < termvec[i+1].Id)
	}
}

func TestWord(t *testing.T) {
	id2word := map[int]string{1: "a", 2: "b", 3: "c"}

	vocab := &Vocabulary{
		id2word: id2word,
	}

	for id, expectedWord := range id2word {
		assert.Equal(t, expectedWord, vocab.Word(id))
	}

	assert.Equal(t, "", vocab.Word(9999))
}

func TestVectorize(t *testing.T) {
	vocab := &Vocabulary{
		word2id:    map[string]int{"a": 1, "b": 2, "c": 3},
		id2word:    map[int]string{1: "a", 2: "b", 3: "c"},
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
