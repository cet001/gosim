package gosim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTFIDF_AddDoc(t *testing.T) {
	c := NewTFIDF()

	docvec1 := []Term{
		{Id: 1, Score: 10},
		{Id: 2, Score: 20},
	}

	docvec2 := []Term{
		{Id: 1, Score: 100},
		{Id: 3, Score: 300},
	}

	// Add 2 document vectors
	c.needsRecalc = false
	c.AddDoc(10000, docvec1)
	c.AddDoc(20000, docvec2)

	assert.Equal(t, 2, len(c.docs))
	assert.True(t, c.needsRecalc)
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
		float64((2*4)+(3*5)), /* expected */
		dot(
			[]Term{{100, 2}, {101, 3}},
			[]Term{{100, 4}, {101, 5}},
		),
	)

	assert.Equal(t,
		float64((2*4)+(3*5)+(7*0)+(0*8)), /* expected */
		dot(
			[]Term{{100, 2}, {101, 3}, {102, 7}},
			[]Term{{100, 4}, {101, 5}, {103, 8}},
		),
	)

	assert.Equal(t,
		float64((-2*0)+(0*3)+(2*-4)), /* expected */
		dot(
			[]Term{{100, -2}, {101, 0}, {102, 2}},
			[]Term{{100, 0}, {101, 3}, {102, -4}},
		),
	)
}
