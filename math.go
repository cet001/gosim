package gosim

import (
	"math"
)

type Term struct {
	Id    int
	Value float64
}

// Sorts terms by increasing Id.
type byTermId []Term

func (a byTermId) Len() int           { return len(a) }
func (a byTermId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTermId) Less(i, j int) bool { return a[i].Id < a[j].Id }

// Sorts terms by increasing Value.
type byTermValue []Term

func (a byTermValue) Len() int           { return len(a) }
func (a byTermValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTermValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

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

// Calculates the Euclidean norm (a.k.a. L2-norm) of the specified vector.
func Norm(vec SparseVector) float64 {
	sumOfSquares := 0.0
	for i := 0; i < len(vec); i++ {
		term := &vec[i]
		sumOfSquares += (term.Value * term.Value)
	}

	return math.Sqrt(sumOfSquares)
}
