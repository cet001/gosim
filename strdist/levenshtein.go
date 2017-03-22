package strdist

import (
	"github.com/cet001/gosim/math"
)

// Calculates theLevenshtein distance between 2 strings.
//
// See https://en.wikipedia.org/wiki/Levenshtein_distance
type Levenshtein struct {
	workspace []int
}

// Creates a new Levenshtein calculator with a 1K working buffer (i.e. it will
// only handle strings whose length is less than 1024).
func NewLevenshtein() *Levenshtein {
	workspaceSize := 1024
	return &Levenshtein{
		workspace: make([]int, workspaceSize),
	}
}

func (me *Levenshtein) Dist(a, b string) int {
	// This implementation was copied and modified from:
	//     https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Levenshtein_distance#Go
	f := me.workspace[0 : len(b)+1]

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := math.Min(f[j], f[j-1]) + 1 // delete & insert
			if cb != ca {
				mn = math.Min(mn, fj1+1) // change
			} else {
				mn = math.Min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}
