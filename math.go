package gosim

import (
	"math"
)

type Term struct {
	Id    int
	Value float64
}

// Sorts terms by increasing Id.
type ByTermId []Term

func (a ByTermId) Len() int           { return len(a) }
func (a ByTermId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTermId) Less(i, j int) bool { return a[i].Id < a[j].Id }

// Sorts terms by decreasing Value.
type ByTermValueDesc []Term

func (a ByTermValueDesc) Len() int           { return len(a) }
func (a ByTermValueDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTermValueDesc) Less(i, j int) bool { return a[i].Value > a[j].Value }

// Represents a sparse vector, where most of the elements are typically empty.
type SparseVector []Term

// Calculates the dot product of vectors v1 and v2.
func Dot(v1, v2 SparseVector) float64 {
	var dotProduct float64 = 0.0
	lenV1, lenV2 := len(v1), len(v2)
	idx1, idx2 := 0, 0

	for {
		if idx1 == lenV1 || idx2 == lenV2 {
			break
		}

		term1, term2 := &v1[idx1], &v2[idx2]

		if term1.Id < term2.Id {
			idx1++
		} else if term2.Id < term1.Id {
			idx2++
		} else {
			dotProduct += (term1.Value * term2.Value)
			idx1++
			idx2++
		}
	}

	return dotProduct
}

// Calculates the Euclidean norm (a.k.a. L2-Norm) of the specified vector.
func Norm(vec SparseVector) float64 {
	sumOfSquares := 0.0
	for i := 0; i < len(vec); i++ {
		term := &vec[i]
		sumOfSquares += (term.Value * term.Value)
	}

	return math.Sqrt(sumOfSquares)
}

// Calculates a weighted mean for the specified values and associated weights.
// This function assumes that:
//    - x and w are are the same length
//    - all values in x and w are non-negative
//    - the sum of the weights is > 0
func WeightedMean(x, w []float64) float64 {
	sumOfWeightedValues := 0.0
	sumOfWeights := 0.0
	for i, xVal := range x {
		sumOfWeightedValues += (xVal * w[i])
		sumOfWeights += w[i]
	}

	return sumOfWeightedValues / sumOfWeights
}

// Hashes a string into an int.  This operation is useful if you want to convert
// a weighted string vector into a SparseVector. I.e. each string s[i] in a
// weighted vector V=[s0w0, s1w1, ...] can be converted into an int.
func Hash(s string) int {
	h := 1125899906842597 // prime
	len := len(s)

	for i := 0; i < len; i++ {
		h = 31*h + int(s[i])
	}

	return h
}

// Similar to the Unix 'uniq' command, this function removes all dupicates from
// a sorted array of int values.
func Uniq(sortedValues []int) []int {
	if sortedValues == nil {
		return []int{}
	}

	if len(sortedValues) <= 1 {
		return sortedValues
	}

	uniqueValues := make([]int, 0, len(sortedValues))
	uniqueValues = append(uniqueValues, sortedValues[0])

	for i := 1; i < len(sortedValues); i++ {
		if sortedValues[i] != sortedValues[i-1] {
			uniqueValues = append(uniqueValues, sortedValues[i])
		}
	}

	return uniqueValues
}

// Returns the intersection of 2 sorted sets.
// Unpredicatable results ensue if a or b contain duplicate elements or are not
// in ascending sorted order.
func Intersect(a, b []int) []int {
	lenA, lenB := len(a), len(b)
	intersection := make([]int, 0, min(lenA, lenB))

	idx1, idx2 := 0, 0
	for {
		if idx1 == lenA || idx2 == lenB {
			break
		}

		aVal, bVal := a[idx1], b[idx2]

		if aVal < bVal {
			idx1++
		} else if bVal < aVal {
			idx2++
		} else {
			intersection = append(intersection, aVal)
			idx1++
			idx2++
		}
	}

	return intersection
}

// Binary merge of sorted sets a and b.
// Unpredicatable results ensue if a or b contain duplicate elements or are not
// in ascending sorted order.
func Union(a, b []int) []int {
	lenA, lenB := len(a), len(b)
	union := make([]int, 0, max(lenA, lenB))

	idx1, idx2 := 0, 0
	for {
		if idx1 == lenA {
			union = append(union, b[idx2:]...)
			break
		} else if idx2 == lenB {
			union = append(union, a[idx1:]...)
			break
		}

		aVal, bVal := a[idx1], b[idx2]

		if aVal < bVal {
			union = append(union, aVal)
			idx1++
		} else if bVal < aVal {
			union = append(union, bVal)
			idx2++
		} else {
			union = append(union, aVal)
			idx1++
			idx2++
		}
	}

	return union
}
