package gosim

import (
	"encoding/gob"
	"github.com/cet001/gosim/math"
	"os"
	"sort"
)

// Vectorize() converts an array of words (terms like "car" or "john smith")
// into a term frequency feature vector where each term is assigned a unique
// integer Id and term frequency.
//
// If 'updateDict' is true, then new encountered terms will be added to the
// underlying Dictionary.
//
// Must return the array of terms SORTED by increasing Term.Id.
// type Vectorize func(words []string, updateDict bool) []math.Term

// Manages the mapping between words and their corresponding integer IDs.
type Dictionary struct {
	word2id    map[string]int
	id2word    map[int]string
	nextTermId int
}

func NewDictionary() *Dictionary {
	const initialCapacity = 1000000
	return &Dictionary{
		word2id:    make(map[string]int, initialCapacity),
		id2word:    make(map[int]string, initialCapacity),
		nextTermId: 1,
	}
}

// Returns the number of words in this dictionary.
func (me *Dictionary) Size() int {
	return len(me.word2id)
}

// Returns the source word (token) corresponding to the the specified term Id.
func (me *Dictionary) Word(termId int) string {
	word, _ := me.id2word[termId]
	return word
}

// Removes the specified terms from this dictionary.
// Returns the number of terms that were removed.
func (me *Dictionary) Remove(terms []math.Term) int {
	numTermsRemoved := 0

	for _, term := range terms {
		word, found := me.id2word[term.Id]
		if found {
			delete(me.id2word, term.Id)
			delete(me.word2id, word)
			numTermsRemoved++
		}
	}

	return numTermsRemoved
}

// Vectorize() converts an array of words (terms like "car" or "john smith")
// into a term frequency feature vector where each term is assigned a unique
// integer Id and term frequency.
//
// If 'updateDict' is true, then new encountered terms will be added to the
// underlying Dictionary.
//
// Returns the term freqency feature vector in sorted order by increasing Term.Id.
func (me *Dictionary) Vectorize(words []string, updateDict bool) math.SparseVector {
	word2freq := make(map[string]int, len(words))
	for _, word := range words {
		word2freq[word]++
	}

	terms := make([]math.Term, 0, len(word2freq))

	if updateDict {
		for word, freq := range word2freq {
			termId, found := me.word2id[word]
			if !found {
				termId = me.nextTermId
				me.word2id[word] = termId
				me.id2word[termId] = word
				me.nextTermId++
			}
			terms = append(terms, math.Term{Id: termId, Value: float64(freq)})
		}
	} else {
		for word, freq := range word2freq {
			termId, found := me.word2id[word]
			if found {
				terms = append(terms, math.Term{Id: termId, Value: float64(freq)})
			}
		}
	}

	sort.Sort(math.ByTermId(terms))
	return terms
}

// Saves the specified Dictionary object to a binary file.
func SaveDictionary(d *Dictionary, filePath string) error {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(len(d.word2id))
		encoder.Encode(d.nextTermId)
		encoder.Encode(d.word2id)
	}

	return err
}

// Loads a Dictionary from the specified binary file.
func LoadDictionary(filePath string) (*Dictionary, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(file)
	d := &Dictionary{}
	var dictSize int

	decodeFuncs := []func() error{
		func() error {
			return decoder.Decode(&dictSize)
		},
		func() error {
			return decoder.Decode(&d.nextTermId)
		},
		func() error {
			return decoder.Decode(&d.word2id)
		},
	}

	for _, decode := range decodeFuncs {
		err := decode()
		if err != nil {
			return nil, err
		}
	}

	// Build the reverse lookup
	d.id2word = make(map[int]string, dictSize)
	for k, v := range d.word2id {
		d.id2word[v] = k
	}

	return d, nil

}
