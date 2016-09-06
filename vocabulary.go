package gosim

import (
	"encoding/gob"
	"os"
	"sort"
)

// Converts an array of words (i.e. tokens like "car", "john smith") into a
// term frequency feature vector where each term is assigned a unique integer
// Id and and a term frequency.
//
// If 'updateVocab' is true, then new encountered terms will be added to the
// underlying Vocablulary.
//
// Must return the array of terms SORTED by increasing Term.Id.
//
type Vectorize func(words []string, updateVocab bool) []Term

// Manages the vocabulary (i.e. the set of distinct terms) for a given corpus.
type Vocabulary struct {
	word2id    map[string]int
	id2word    map[int]string
	nextTermId int
}

func NewVocabulary() *Vocabulary {
	const initialCapacity = 1000000
	return &Vocabulary{
		word2id:    make(map[string]int, initialCapacity),
		id2word:    make(map[int]string, initialCapacity),
		nextTermId: 1,
	}
}

// Returns the number of terms in this vocabulary.
func (me *Vocabulary) Size() int {
	return len(me.word2id)
}

// Returns the source word (token) corresponding to the the specified term Id.
func (me *Vocabulary) Word(termId int) string {
	word, _ := me.id2word[termId]
	return word
}

// Removes the specified terms from this vocabulary.
// Returns the number of terms that were removed.
func (me *Vocabulary) Remove(terms []Term) int {
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

// See the Vectorize() function def at the top of this file.
func (me *Vocabulary) Vectorize(words []string, updateVocab bool) []Term {
	word2freq := make(map[string]int, len(words))
	for _, word := range words {
		word2freq[word]++
	}

	terms := make([]Term, 0, len(word2freq))

	if updateVocab {
		for word, freq := range word2freq {
			termId, found := me.word2id[word]
			if !found {
				termId = me.nextTermId
				me.word2id[word] = termId
				me.id2word[termId] = word
				me.nextTermId++
			}
			terms = append(terms, Term{Id: termId, Value: float64(freq)})
		}
	} else {
		for word, freq := range word2freq {
			termId, found := me.word2id[word]
			if found {
				terms = append(terms, Term{Id: termId, Value: float64(freq)})
			}
		}
	}

	sort.Sort(byTermId(terms))
	return terms
}

// Saves the specified Vocabulary object to a binary file.
func SaveVocabulary(vocab *Vocabulary, filePath string) error {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(len(vocab.word2id))
		encoder.Encode(vocab.nextTermId)
		encoder.Encode(vocab.word2id)
	}

	return err
}

// Loads a Vocabulary from the specified binary file.
func LoadVocabulary(filePath string) (*Vocabulary, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(file)
	vocab := &Vocabulary{}
	var vocabSize int

	decodeFuncs := []func() error{
		func() error {
			return decoder.Decode(&vocabSize)
		},
		func() error {
			return decoder.Decode(&vocab.nextTermId)
		},
		func() error {
			return decoder.Decode(&vocab.word2id)
		},
	}

	for _, decode := range decodeFuncs {
		err := decode()
		if err != nil {
			return nil, err
		}
	}

	// Build the reverse lookup
	vocab.id2word = make(map[int]string, vocabSize)
	for k, v := range vocab.word2id {
		vocab.id2word[v] = k
	}

	return vocab, nil

}
