package gosim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
)

func TestByTermValueDesc(t *testing.T) {
	terms := []Term{{1, 0.01}, {2, 0.02}, {3, 0.03}}
	sort.Sort(ByTermValueDesc(terms))
	assert.Equal(t, []Term{{3, 0.03}, {2, 0.02}, {1, 0.01}}, terms)
}

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

func TestIntersect(t *testing.T) {
	assert.Equal(t, []int{}, Intersect([]int{}, []int{}))
	assert.Equal(t, []int{}, Intersect(nil, nil))
	assert.Equal(t, []int{}, Intersect([]int{}, []int{1, 2, 3}))
	assert.Equal(t, []int{}, Intersect([]int{1, 2, 3}, []int{}))
	assert.Equal(t, []int{}, Intersect([]int{1, 2}, []int{3, 4}))

	assert.Equal(t, []int{1}, Intersect([]int{1, 2, 3}, []int{1, 4}))
	assert.Equal(t, []int{1, 3}, Intersect([]int{1, 2, 3}, []int{1, 3, 4}))
	assert.Equal(t, []int{2, 3}, Intersect([]int{1, 2, 3}, []int{2, 3, 4}))
	assert.Equal(t, []int{2, 3}, Intersect([]int{2, 3, 4}, []int{1, 2, 3}))

	assert.Equal(t, []int{1, 2, 3}, Intersect([]int{1, 2, 3}, []int{1, 2, 3}))
}

func BenchmarkIntersect_Small(b *testing.B) {
	setA, setB := make([]int, 1000), make([]int, 1000)
	fmt.Printf("len(a)=%v, len(b)=%v\n", len(setA), len(setB))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Intersect(
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			[]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25},
		)
	}
}

func BenchmarkIntersect_Big(b *testing.B) {
	setA, setB := make([]int, 1000), make([]int, 1000)
	for i := 0; i < len(setA); i++ {
		setA[i] = i
		setB[i] = i + (len(setA) / 3)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Intersect(setA, setB)
	}
}
