package gosim

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestDot(t *testing.T) {
	assert.Equal(t,
		float64((2*4)+(3*5)), // expected
		Dot(
			[]Term{{100, 2}, {101, 3}},
			[]Term{{100, 4}, {101, 5}},
		),
	)

	assert.Equal(t,
		float64((2*4)+(3*5)+(7*0)+(0*8)), // expected
		Dot(
			[]Term{{100, 2}, {101, 3}, {102, 7}},
			[]Term{{100, 4}, {101, 5}, {103, 8}},
		),
	)

	assert.Equal(t,
		float64((-2*0)+(0*3)+(2*-4)), // expected
		Dot(
			[]Term{{100, -2}, {101, 0}, {102, 2}},
			[]Term{{100, 0}, {101, 3}, {102, -4}},
		),
	)

	assert.Equal(t,
		float64(0), // expected
		Dot(
			[]Term{},
			[]Term{{100, 1}, {101, 2}, {102, 3}},
		),
	)
}

func BenchmarkDot(b *testing.B) {
	const vecSize = 10000
	rnd := rand.New(rand.NewSource(99))

	makeRandomVector := func(size int) SparseVector {
		v := make(SparseVector, 0, size)
		for i := 0; i < size; i++ {
			v = append(v, Term{Id: i, Value: rnd.Float64()})
		}
		return v
	}

	v1, v2 := makeRandomVector(vecSize), makeRandomVector(vecSize)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Dot(v1, v2)
	}
}

func TestNorm(t *testing.T) {
	// sqrt(2^2 + 3^2 + 6^2) = 7
	assert.Equal(t, 7.0, Norm([]Term{{100, 2}, {101, 3}, {102, 6}}))

	// sqrt(0^2 + 0^2) = 0
	assert.Equal(t, 0.0, Norm([]Term{{100, 0}, {101, 0}}))

	// sqrt(5^2 + 0^2) = 5
	assert.Equal(t, 5.0, Norm([]Term{{100, 5}, {101, 0}}))
}

func TestWeightedMean(t *testing.T) {
	x := []float64{10.0, 20.0, 30.0}
	w := []float64{0.20, 0.30, 0.50}
	assert.Equal(t, ((10.0 * 0.20) + (20.0 * 0.30) + (30.0*0.50)/(0.20+0.30+0.50)), WeightedMean(x, w))
}

// This is just a sanity-check.
func TestHash(t *testing.T) {
	values := []string{
		"", "a", "b", "c", "A", "B", "C", "cat", "CAT",
		"aaaaaaaaaaaaaaaa", "???????????????????????",
		"1", " 1", "  1",
	}

	uniqueHashValues := map[int]bool{}
	for _, value := range values {
		uniqueHashValues[Hash(value)] = true
	}

	assert.Equal(t, len(values), len(uniqueHashValues))
}
