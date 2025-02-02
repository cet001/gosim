package gosim

import (
	"fmt"
	"github.com/cet001/mathext/vectors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sort"
	"testing"
)

func ExampleDictionary() {
	d := NewDictionary()
	doc := []string{"three", "one", "two", "three", "two", "three"}
	termvec := d.VectorizeAndUpdate(doc)

	fmt.Printf("Dictionary has %v distinct terms\n", d.Size())

	sort.Sort(vectors.ByElementValueDesc(termvec))
	for _, term := range termvec {
		fmt.Printf("'%v' term frequency is %v\n", d.Word(term.Id), term.Value)
	}
	// Output:
	// Dictionary has 3 distinct terms
	// 'three' term frequency is 0.5
	// 'two' term frequency is 0.3333333333333333
	// 'one' term frequency is 0.16666666666666666
}

func TestDictionary_BasicUsage(t *testing.T) {
	d := NewDictionary()
	assert.Equal(t, 0, d.Size())

	termvec := d.VectorizeAndUpdate([]string{"b", "a", "c", "b", "a", "a"})

	// Verify terms are ordered by Id
	for i := 0; i < len(termvec)-1; i++ {
		assert.True(t, termvec[i].Id < termvec[i+1].Id)
	}
}

func TestDictionary_Word(t *testing.T) {
	id2word := map[int]string{1: "a", 2: "b", 3: "c"}

	d := &Dictionary{
		id2word: id2word,
	}

	for id, expectedWord := range id2word {
		assert.Equal(t, expectedWord, d.Word(id))
	}

	assert.Equal(t, "", d.Word(9999))
}

func TestDictionary_Vectorize(t *testing.T) {
	d := &Dictionary{
		word2id:    map[string]int{"a": 1, "b": 2, "c": 3},
		id2word:    map[int]string{1: "a", 2: "b", 3: "c"},
		nextTermId: 4,
	}

	vec := d.Vectorize([]string{"c", "a", "a", "Z", "Z", "Z"})
	assert.Equal(t, vectors.SparseVector{{Id: 1, Value: 2.0 / 6.0}, {Id: 3, Value: 1.0 / 6.0}}, vec)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, d.word2id)
}

func TestDictionary_VectorizeAndUpdate(t *testing.T) {
	d := &Dictionary{
		word2id:    map[string]int{"a": 1, "b": 2, "c": 3},
		id2word:    map[int]string{1: "a", 2: "b", 3: "c"},
		nextTermId: 4,
	}

	vec := d.VectorizeAndUpdate([]string{"c", "a", "a", "Z", "Z", "Z"})
	assert.Equal(t, vectors.SparseVector{{Id: 1, Value: 2.0 / 6.0}, {Id: 3, Value: 1.0 / 6.0}, {Id: 4, Value: 3.0 / 6.0}}, vec)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3, "Z": 4}, d.word2id)
}

func TestDictionary_Remove(t *testing.T) {
	d := &Dictionary{
		word2id:    map[string]int{"a": 1, "b": 2, "c": 3},
		id2word:    map[int]string{1: "a", 2: "b", 3: "c"},
		nextTermId: 4,
	}

	numTermsRemoved := d.Remove([]vectors.Element{
		{Id: 1, Value: 100},
		{Id: 3, Value: 300},
		{Id: 4, Value: 400},
	})

	assert.Equal(t, 2, numTermsRemoved)
	assert.Equal(t, map[string]int{"b": 2}, d.word2id)
	assert.Equal(t, map[int]string{2: "b"}, d.id2word)
}

func TestDictionary_SaveAndLoad(t *testing.T) {
	d := &Dictionary{
		word2id:    map[string]int{"a": 1, "b": 2},
		id2word:    map[int]string{1: "a", 2: "b"},
		nextTermId: 3,
	}

	f, _ := ioutil.TempFile("/tmp", "gosim_test_")
	dictFilePath := f.Name()
	defer os.Remove(dictFilePath)

	err := SaveDictionary(d, dictFilePath)
	assert.Nil(t, err)

	d, err = LoadDictionary(dictFilePath)
	if assert.Nil(t, err) {
		assert.Equal(t, map[string]int{"a": 1, "b": 2}, d.word2id)
		assert.Equal(t, map[int]string{1: "a", 2: "b"}, d.id2word)
		assert.Equal(t, 3, d.nextTermId)
	}
}
