package gosim

// Calculates the Jaro-Winkler string distance (similarity) score.
// Based on the Wikipedia description: https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance.
type JaroWinkerDist struct {
	s1matches   []bool
	s2matches   []bool
	prefixScale float64
}

func NewJaroWinklerDist() *JaroWinkerDist {
	maxStringLen := 1024 * 4 // size of working space

	return &JaroWinkerDist{
		s1matches:   make([]bool, maxStringLen),
		s2matches:   make([]bool, maxStringLen),
		prefixScale: 0.1,
	}
}

// Returns the  Jaro-Winkler distance score for s1 and s2.
// WARNING: This call is NOT threadsafe!
func (me *JaroWinkerDist) Calc(s1, s2 []byte) float64 {
	lenS1, lenS2 := len(s1), len(s2)
	maxMatchDist := (max(lenS1, lenS2) / 2) - 1
	s1matches, s2matches := me.s1matches, me.s2matches

	// Clear the working space
	for i := 0; i < max(lenS1, lenS2); i++ {
		s1matches[i] = false
		s2matches[i] = false
	}

	// Count the matches and track which characters from each string matched.
	m := float64(0)
	for i, ch1 := range s1 {
		left, right := max(0, i-maxMatchDist), min(lenS2, i+maxMatchDist+1)
		for j := left; j < right; j++ {
			ch2 := s2[j]
			if ch1 == ch2 && !s2matches[j] {
				m++
				s1matches[i] = true
				s2matches[j] = true
			}
		}
	}

	if m == 0 { // no matches
		return 0.0
	}

	// Calculate the number of full transpositions
	halfTranspositions := 0
	j := 0
	for i := 0; i < lenS1; i++ {
		if s1matches[i] {
			for !s2matches[j] {
				j++
			}
			if s1[i] != s2[j] {
				halfTranspositions++
			}
			j++
		}
	}
	t := float64(halfTranspositions / 2)

	jaroDist := ((m / float64(lenS1)) + (m / float64(lenS2)) + ((m - t) / m)) / 3.0

	// Calculate the length of the largest common prefix of s1 and s2
	p := 0
	for i := 0; i < lenS1; i++ {
		if i < lenS2 {
			if s1[i] != s2[i] {
				break
			}
			p++
		}
	}
	jaroWinklerDist := jaroDist + float64(p)*me.prefixScale*(1.0-jaroDist)
	return jaroWinklerDist
}

// Returns the lesser of a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Returns the greater of a and b.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
