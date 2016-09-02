package gosim

import (
	"sort"
)

// Converts an array of words (i.e. tokens like "car", "john smith") into a
// "term frequency" feature vector where each term is assigned a unique integer
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

func (me *Vocabulary) Size() int {
	return len(me.word2id)
}

func (me *Vocabulary) Word(termId int) string {
	word, _ := me.id2word[termId]
	return word
}

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
