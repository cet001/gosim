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

func TestNorm(t *testing.T) {
	// sqrt(2^2 + 3^2 + 6^2) = 7
	assert.Equal(t, 7.0, norm([]Term{{100, 2}, {101, 3}, {102, 6}}))

	// sqrt(0^2 + 0^2) = 0
	assert.Equal(t, 0.0, norm([]Term{{100, 0}, {101, 0}}))

	// sqrt(5^2 + 0^2) = 5
	assert.Equal(t, 5.0, norm([]Term{{100, 5}, {101, 0}}))
}

func TestDot(t *testing.T) {
	assert.Equal(t,
		float64((2*4)+(3*5)), // expected
		dot(
			[]Term{{100, 2}, {101, 3}},
			[]Term{{100, 4}, {101, 5}},
		),
	)

	assert.Equal(t,
		float64((2*4)+(3*5)+(7*0)+(0*8)), // expected
		dot(
			[]Term{{100, 2}, {101, 3}, {102, 7}},
			[]Term{{100, 4}, {101, 5}, {103, 8}},
		),
	)

	assert.Equal(t,
		float64((-2*0)+(0*3)+(2*-4)), // expected
		dot(
			[]Term{{100, -2}, {101, 0}, {102, 2}},
			[]Term{{100, 0}, {101, 3}, {102, -4}},
		),
	)

	assert.Equal(t,
		float64(0), // expected
		dot(
			[]Term{},
			[]Term{{100, 1}, {101, 2}, {102, 3}},
		),
	)
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

	// Create a filter that specifies which term Ids are to be kept.
	filter := map[int]int{10: 999, 30: 999, 50: 999}

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
