package gosim

import (
	"fmt"
	"math/rand"
	"time"
)

type LDATerm struct {
	Id    int // this term's unique Id within the corpus vocabulary
	Topic int // the index of the topic of which this term is a member
}

type LDADocument struct {
	// The number of terms contained by topic[i] within this document.
	NumTermsInTopic []int
}

// s1 := rand.NewSource(time.Now().UnixNano())
// r1 := rand.New(s1)

// Implementation of the LDA (Latent Dirichlet Allocation) model.
// Based on http://brooksandrew.github.io/simpleblog/articles/latent-dirichlet-allocation-under-the-hood/
type LDAModel struct {
	K int // number of topics

	// Hyperparameter:
	//    Alpha=1 : symmetric dirichlet prior
	//    Alpha>1 : scatters document clusters
	Alpha int

	Eta float64

	// Topic-Word matrix.
	tw [][]int
}

func NewLDAModel() *LDAModel {
	return &LDAModel{
		K:     20,
		Alpha: 1,
		Eta:   0.001,
	}
}

func (me *LDAModel) Train(vocabSize int, nextDoc func() (int, SparseVector, bool)) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	logger.Printf("Initializing topic-word matrix")
	// twMatrix = make([]int, me.K*vocabSize)

	logger.Printf("")
	docIdx, docVector, hasMoreDocs := nextDoc()
	for hasMoreDocs {
		for _, term := range docVector {
			randomTopicIdx := rnd.Intn(me.K)
			fmt.Printf("doc=%v; term=%v: randomTopicIdx=%v", docIdx, term.Id, randomTopicIdx)
			// wt[randomTopicIdx, term.Id] += 1
		}

		docIdx, docVector, hasMoreDocs = nextDoc()
	}
}
