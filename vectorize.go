package gosim

import (
	"sort"
)

type Vectorize func(words []string, updateVocab bool) []Term

type Vocabulary struct {
	word2id    map[string]int
	nextTermId int
}

func NewVocabulary() *Vocabulary {
	return &Vocabulary{
		word2id:    make(map[string]int, 1000000),
		nextTermId: 1,
	}
}

func (me *Vocabulary) Size() int {
	return len(me.word2id)
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
