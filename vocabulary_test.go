package gosim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sort"
	"testing"
)

func ExampleVocabulary() {
	vocab := NewVocabulary()
	doc := []string{"row", "row", "row", "your", "boat"}
	termvec := vocab.Vectorize(doc, true)

	fmt.Printf("Vocabulary has %v distinct terms\n", vocab.Size())

	sort.Sort(ByTermValueDesc(termvec))
	for _, term := range termvec {
		fmt.Printf("'%v' has %v occurences\n", vocab.Word(term.Id), term.Value)
	}
	// Output:
	// Vocabulary has 3 distinct terms
	// 'row' has 3 occurences
	// 'your' has 1 occurences
	// 'boat' has 1 occurences
}

func TestVocabulary_BasicUsage(t *testing.T) {
	vocab := NewVocabulary()
	assert.Equal(t, 0, vocab.Size())

	termvec := vocab.Vectorize([]string{"b", "a", "c", "b", "a", "a"}, true)

	// Verify terms are ordered by Id
	for i := 0; i < len(termvec)-1; i++ {
		assert.True(t, termvec[i].Id < termvec[i+1].Id)
	}
}

func TestVocabulary_Word(t *testing.T) {
	id2word := map[int]string{1: "a", 2: "b", 3: "c"}

	vocab := &Vocabulary{
		id2word: id2word,
	}

	for id, expectedWord := range id2word {
		assert.Equal(t, expectedWord, vocab.Word(id))
	}

	assert.Equal(t, "", vocab.Word(9999))
}

func TestVocabulary_Vectorize(t *testing.T) {
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

func TestVocabulary_Remove(t *testing.T) {
	vocab := &Vocabulary{
		word2id:    map[string]int{"a": 1, "b": 2, "c": 3},
		id2word:    map[int]string{1: "a", 2: "b", 3: "c"},
		nextTermId: 4,
	}

	numTermsRemoved := vocab.Remove([]Term{{Id: 1, Value: 100}, {Id: 3, Value: 300}, {Id: 4, Value: 400}})
	assert.Equal(t, 2, numTermsRemoved)
	assert.Equal(t, map[string]int{"b": 2}, vocab.word2id)
	assert.Equal(t, map[int]string{2: "b"}, vocab.id2word)
}

func TestVocabulary_SaveAndLoad(t *testing.T) {
	vocab := &Vocabulary{
		word2id:    map[string]int{"a": 1, "b": 2},
		id2word:    map[int]string{1: "a", 2: "b"},
		nextTermId: 3,
	}

	f, _ := ioutil.TempFile("/tmp", "gosim_test_")
	vocabFilePath := f.Name()
	defer os.Remove(vocabFilePath)

	err := SaveVocabulary(vocab, vocabFilePath)
	assert.Nil(t, err)

	vocab, err = LoadVocabulary(vocabFilePath)
	assert.Nil(t, err)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, vocab.word2id)
	assert.Equal(t, map[int]string{1: "a", 2: "b"}, vocab.id2word)
	assert.Equal(t, 3, vocab.nextTermId)
}
