package gosim

import (
	"math"
)

// Term represents a token (typically a word) having a unique ID within a document.
type Term struct {
	Id    int
	Value float64
}

// Sorts Term objects by increasing Term.Id.
type ByTermId []Term

func (a ByTermId) Len() int           { return len(a) }
func (a ByTermId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTermId) Less(i, j int) bool { return a[i].Id < a[j].Id }

// Sorts Term objects by decreasing Term.Value.
type ByTermValueDesc []Term

func (a ByTermValueDesc) Len() int           { return len(a) }
func (a ByTermValueDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTermValueDesc) Less(i, j int) bool { return a[i].Value > a[j].Value }

// Represents a sparse vector, where most of the elements are typically empty.
// For example, consider the following vector containing 10 elements:
//
//   v = [9 0 0 2 0 0 0 0 7 0]
//
// Only elements 0, 3, and 8 (base-0) contain non-zero values -- the remaining
// elements are "empty".  The following sparse vector is equivalent to the above
// vector:
//
//   sv := SparseVector{{0, 9}, {3, 2}, {8, 7}}
//
// Each element in the SparseVector is a Term object that specifies the element's
// value (Term.Value) and position (Term.Id) within the vector.  Note that the
// SparseVector declaration above is a shorthand syntax; it can also be declared
// more formally like this:
//
//   sv := SparseVector{
//	   Term{Id: 0, Value: 9},
//	   Term{Id: 3, Value: 2},
//	   Term{Id: 8, Value: 7},
//   }
//
type SparseVector []Term

// Calculates the dot product of two sparse vectors.  Dot() assumes that v1 and
// v2 are in sorted order by Term.Id.
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

// Calculates a weighted mean for the specified values in []x and associated
// weights in []w.
//
// This function assumes that:
//    - x and w are are the same length
//    - all values in x and w are non-negative
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
//
// a and b are the sets to be intersected.
//
// target is an optional slice into which the intersecting elements from a and b
// are appended.  If target is nil, a new []int slice will be created and returned.
//
// WARNING: Unpredicatable results ensue if a or b contain duplicate elements or
// are not in ascending sorted order.
func Intersect(a, b, target []int) []int {
	lenA, lenB := len(a), len(b)

	var intersection []int
	if target == nil {
		intersection = make([]int, 0, min(lenA, lenB))
	} else {
		intersection = target[:0]
	}

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

// Returns the union of 2 sorted sets.
//
// a and b are the sets to be unioned.
//
// target is an optional slice into which the unique  elements from a and b are
// appended.  If this param is nil, a new []int slice will be created.
//
// WARNING: Unpredicatable results ensue if a or b contain duplicate elements or
// are not in ascending sorted order.

// Binary merge of sorted sets a and b.
// Unpredicatable results ensue if a or b contain duplicate elements or are not
// in ascending sorted order.
func Union(a, b, target []int) []int {
	lenA, lenB := len(a), len(b)

	var union []int
	if target == nil {
		union = make([]int, 0, max(lenA, lenB))
	} else {
		union = target[:0]
	}

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
