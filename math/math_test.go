package math

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
)

func TestMin(t *testing.T) {
	assert.Equal(t, 1, Min(1, 2))
	assert.Equal(t, 1, Min(2, 1))

	assert.Equal(t, -2, Min(-1, -2))
	assert.Equal(t, -2, Min(-2, -1))

	assert.Equal(t, -1, Min(-1, 2))
	assert.Equal(t, -1, Min(2, -1))

	assert.Equal(t, 3, Min(3, 3))
	assert.Equal(t, -3, Min(-3, -3))

	assert.Equal(t, 0, Min(99, 0))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 2, Max(1, 2))
	assert.Equal(t, 2, Max(2, 1))

	assert.Equal(t, -1, Max(-1, -2))
	assert.Equal(t, -1, Max(-2, -1))

	assert.Equal(t, 2, Max(-1, 2))
	assert.Equal(t, 2, Max(2, -1))

	assert.Equal(t, 3, Max(3, 3))
	assert.Equal(t, -3, Max(-3, -3))

	assert.Equal(t, 99, Max(99, 0))
}

func TestByTermId(t *testing.T) {
	terms := SparseVector{{3, 0.30}, {1, 0.10}, {2, 0.20}}
	sort.Sort(ByTermId(terms))
	assert.Equal(t, SparseVector{{1, 0.10}, {2, 0.20}, {3, 0.30}}, terms)
}

func TestByTermValueDesc(t *testing.T) {
	terms := SparseVector{{1, 0.01}, {2, 0.02}, {3, 0.03}}
	sort.Sort(ByTermValueDesc(terms))
	assert.Equal(t, SparseVector{{3, 0.03}, {2, 0.02}, {1, 0.01}}, terms)
}

func ExampleDot() {
	v1 := SparseVector{{0, 4}, {2, 2}}
	v2 := SparseVector{{0, 5}, {1, 3}, {2, 7}, {3, 6}}
	fmt.Println(Dot(v1, v2))
	// Output:
	// 34
}

func TestDot(t *testing.T) {
	assert.Equal(t,
		float64((2*4)+(3*5)), // expected
		Dot(
			SparseVector{{100, 2}, {101, 3}},
			SparseVector{{100, 4}, {101, 5}},
		),
	)

	assert.Equal(t,
		float64((2*4)+(3*5)+(7*0)+(0*8)), // expected
		Dot(
			SparseVector{{100, 2}, {101, 3}, {102, 7}},
			SparseVector{{100, 4}, {101, 5}, {103, 8}},
		),
	)

	assert.Equal(t,
		float64((-2*0)+(0*3)+(2*-4)), // expected
		Dot(
			SparseVector{{100, -2}, {101, 0}, {102, 2}},
			SparseVector{{100, 0}, {101, 3}, {102, -4}},
		),
	)

	assert.Equal(t,
		float64(0), // expected
		Dot(
			SparseVector{},
			SparseVector{{100, 1}, {101, 2}, {102, 3}},
		),
	)
}

func BenchmarkDot(b *testing.B) {
	const vecSize = 10000
	rnd := rand.New(rand.NewSource(99))

	// Returns a vector of the specified size that is filled with random values.
	makeRandomVector := func(size int) SparseVector {
		v := make(SparseVector, 0, size)
		for i := 0; i < size; i++ {
			v = append(v, Term{Id: i, Value: rnd.Float64()})
		}
		return v
	}

	// Create 2 vectors, each containing 1,000 random values
	v1, v2 := makeRandomVector(vecSize), makeRandomVector(vecSize)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Dot(v1, v2)
	}
}

func ExampleNorm() {
	euclidianNorm := Norm(SparseVector{{1000, 2}, {2000, 3}, {3000, 6}})
	fmt.Println(euclidianNorm)
	// Output:
	// 7
}

func TestNorm(t *testing.T) {
	// sqrt(2^2 + 3^2 + 6^2) = 7
	assert.Equal(t, 7.0, Norm(SparseVector{{100, 2}, {101, 3}, {102, 6}}))

	// sqrt(0^2 + 0^2) = 0
	assert.Equal(t, 0.0, Norm(SparseVector{{100, 0}, {101, 0}}))

	// sqrt(5^2 + 0^2) = 5
	assert.Equal(t, 5.0, Norm(SparseVector{{100, 5}, {101, 0}}))
}

func TestWeightedMean(t *testing.T) {
	x := []float64{10.0, 20.0, 30.0}
	w := []float64{0.20, 0.30, 0.50}
	assert.Equal(t, ((10.0 * 0.20) + (20.0 * 0.30) + (30.0*0.50)/(0.20+0.30+0.50)), WeightedMean(x, w))
}

func ExampleHash() {
	fmt.Println(Hash("john"))
	fmt.Println(Hash("12345678"))
	fmt.Println(Hash("XXX_YYY_ZZZ"))
	fmt.Println(Hash("xxx_yyy_zzz"))
	// Output:
	// 6774539739450401392
	// -4898812128727250071
	// -8286756815414078424
	// -8259655320462518136
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

func ExampleUniq() {
	fmt.Println(Uniq([]int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}))
	// Output:
	// [1 2 3 4]
}

func TestUniq(t *testing.T) {
	assert.Equal(t, []int{}, Uniq(nil))
	assert.Equal(t, []int{}, Uniq([]int{}))
	assert.Equal(t, []int{1}, Uniq([]int{1}))
	assert.Equal(t, []int{1, 3, 5}, Uniq([]int{1, 3, 5}))

	assert.Equal(t, []int{1}, Uniq([]int{1, 1, 1}))
	assert.Equal(t, []int{1, 5}, Uniq([]int{1, 1, 5}))
	assert.Equal(t, []int{1, 5}, Uniq([]int{1, 5, 5}))
	assert.Equal(t, []int{1, 3, 5}, Uniq([]int{1, 3, 3, 3, 5}))
	assert.Equal(t, []int{1, 3, 5}, Uniq([]int{1, 1, 3, 3, 3, 5, 5, 5, 5, 5, 5}))
}

func ExampleIntersect() {
	a := []int{1, 3, 5}
	b := []int{2, 3, 4, 5}
	fmt.Println(Intersect(a, b, nil))
	// Output:
	// [3 5]
}

func TestIntersect(t *testing.T) {
	type TestSet struct {
		a, b, expected []int
	}

	testSets := []TestSet{
		{
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			a:        nil,
			b:        nil,
			expected: []int{},
		},
		{
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: []int{},
		},
		{
			a:        nil,
			b:        []int{1, 2, 3},
			expected: []int{},
		},
		{
			a:        []int{1, 2, 3},
			b:        []int{},
			expected: []int{},
		},
		{
			a:        []int{1, 2},
			b:        []int{3, 4},
			expected: []int{},
		},
		{
			a:        []int{1, 2, 3},
			b:        []int{1, 4},
			expected: []int{1},
		},
		{
			a:        []int{1, 2, 3},
			b:        []int{1, 3, 4},
			expected: []int{1, 3},
		},
		{
			a:        []int{1, 2, 3},
			b:        []int{2, 3, 4},
			expected: []int{2, 3},
		},
		{
			a:        []int{2, 3, 4},
			b:        []int{1, 2, 3},
			expected: []int{2, 3},
		},
		{
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, testSet := range testSets {
		assert.Equal(t, testSet.expected, Intersect(testSet.a, testSet.b, nil))
	}

	intersection := make([]int, 1000)
	for _, testSet := range testSets {
		assert.Equal(t, testSet.expected, Intersect(testSet.a, testSet.b, intersection))
	}
}

func ExampleUnion() {
	a := []int{1, 3}
	b := []int{2, 3, 4}
	fmt.Println(Union(a, b, nil))
	// Output:
	// [1 2 3 4]
}

func TestUnion(t *testing.T) {
	type TestSet struct {
		a, b, expected []int
	}

	testSets := []TestSet{
		{
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			a:        []int{1, 2},
			b:        []int{},
			expected: []int{1, 2},
		},
		{
			a:        []int{1, 2},
			b:        nil,
			expected: []int{1, 2},
		},
		{
			a:        []int{},
			b:        []int{1, 2},
			expected: []int{1, 2},
		},
		{
			a:        []int{1, 2},
			b:        []int{1, 2},
			expected: []int{1, 2},
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []int{1, 2},
			expected: []int{1, 2, 3, 4},
		},
		{
			a:        []int{1, 2},
			b:        []int{1, 2, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []int{2, 3},
			expected: []int{1, 2, 3, 4},
		},
		{
			a:        []int{1, 3},
			b:        []int{2, 4},
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, testSet := range testSets {
		assert.Equal(t, testSet.expected, Union(testSet.a, testSet.b, nil))
	}

	union := make([]int, 1000)
	for _, testSet := range testSets {
		assert.Equal(t, testSet.expected, Union(testSet.a, testSet.b, union))
	}
}

func BenchmarkIntersect_Small(b *testing.B) {
	workspace := make([]int, 1000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Intersect(
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			[]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25},
			workspace,
		)
	}
}

func BenchmarkIntersect_Big(b *testing.B) {
	size := 1000
	setA, setB := make([]int, size), make([]int, size)
	for i := 0; i < len(setA); i++ {
		setA[i] = i
		setB[i] = i + (len(setA) / 3)
	}

	workspace := make([]int, size)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Intersect(setA, setB, workspace)
	}
}
