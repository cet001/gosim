package gosim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTFIDF_AddDoc(t *testing.T) {
	c := NewTFIDF()

	docvec1 := []Term{
		{Id: 1, Value: 10},
		{Id: 2, Value: 20},
	}

	docvec2 := []Term{
		{Id: 1, Value: 100},
		{Id: 3, Value: 300},
	}

	// Add 2 document vectors
	c.needsRecalc = false
	c.AddDoc(10000, docvec1)
	c.AddDoc(20000, docvec2)

	assert.Equal(t, 2, len(c.docs))
	assert.True(t, c.needsRecalc)
}

func TestCalcDocFrequencies(t *testing.T) {
	docs := []Document{
		{
			Id: 100,
			tf: SparseVector{{10, 1}, {20, 2}, {30, 3}},
		},
		{
			Id: 200,
			tf: SparseVector{{20, 200}, {30, 300}},
		},
		{
			Id: 300,
			tf: SparseVector{{30, 300}},
		},
	}

	assert.Equal(t, map[int]int{10: 1, 20: 2, 30: 3}, calcDocFrequencies(docs))
}

func TestCalcTFIDF(t *testing.T) {
	tf := SparseVector{{10, 10}, {40, 40}, {50, 50}}
	idf := SparseHashVector{10: 0.1, 20: 0.2, 30: 0.3, 40: 0.4, 50: 0.5}
	assert.Equal(t, SparseVector{{10, (10 * 0.1)}, {40, (40 * 0.4)}, {50, (50 * 0.5)}}, calcTFIDF(tf, idf))
}

func TestRemoveUnimportantTerms(t *testing.T) {
	docFreqs := map[int]int{1: 1, 2: 2, 3: 10, 4: 20, 5: 30}
	removedTerms := removeUnimportantTerms(docFreqs, 100)
	assert.Equal(t, map[int]int{3: 10, 4: 20}, docFreqs)
	assert.Equal(t, 3, len(removedTerms))
}

func TestFilterDocVectors(t *testing.T) {
	docs := []Document{
		{
			Id: 100,
			tf: SparseVector{{10, 1}, {20, 2}, {30, 3}},
		},
		{
			Id: 200,
			tf: SparseVector{{20, 200}, {30, 300}, {40, 400}},
		},
	}

	// Create a filter that specifies which term Ids are to be kept
	someValue := 999
	filter := map[int]int{10: someValue, 30: someValue, 50: someValue}

	filterDocVectors(docs, filter)

	assert.Equal(t,
		[]Document{
			{
				Id: 100,
				tf: SparseVector{{10, 1}, {30, 3}},
			},
			{
				Id: 200,
				tf: SparseVector{{30, 300}},
			},
		},
		docs,
	)
}
