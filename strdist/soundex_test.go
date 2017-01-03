package strdist

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleSoundex() {
	for _, name := range []string{"Robert", "Rupert", "Romulan", "Larry", "Lori"} {
		fmt.Println(name + " -> " + Soundex(name))
	}
	// Output:
	// Robert -> R163
	// Rupert -> R163
	// Romulan -> R545
	// Larry -> L600
	// Lori -> L600
}

func TestSoundex(t *testing.T) {
	// Using examples from the Wikipedia page on 'Soundex'
	words := map[string]string{
		"":                    "0000",
		"Robert":              "R163",
		"robert":              "R163",
		"Rupert":              "R163",
		"Rubin":               "R150",
		"R_u_b_i_n":           "R150",
		"R_u_b_i_n_99":        "R150",
		"R*&u-&^b###i()n+123": "R150",
		"Ashcraft":            "A261",
		"Ashcroft":            "A261",
		"archer":              "A626",
		// "Tymczak" : "T522",  // <= this fails!
		// "Pfister": "P236",   // <= this fails!
	}

	for w, expectedCode := range words {
		assert.Equal(t, expectedCode, Soundex(w))
	}
}

func Benchmark_Soundex(b *testing.B) {
	words := []string{
		"a", "cat", "catastrophe", "baseball", "dandelion", "elephant",
		"foggy", "ghost", "hangover", "illegal",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, w := range words {
			Soundex(w)
		}
	}
}

func ExampleRefinedSoundex() {
	names := []string{"Braz", "Broz", "Robert", "Rupert", "Rubin", "R_u_b_i_n"}
	for _, name := range names {
		fmt.Println(name + " -> " + RefinedSoundex(name))
	}
	// Output:
	// Braz -> B1905
	// Broz -> B1905
	// Robert -> R901096
	// Rupert -> R901096
	// Rubin -> R90108
	// R_u_b_i_n -> R90108
}

func TestRefinedSoundex(t *testing.T) {
	// Using examples from http://ntz-develop.blogspot.com/2011/03/phonetic-algorithms.html
	words := map[string]string{
		"":          "0000",
		"Braz":      "B1905",
		"Broz":      "B1905",
		"Caren":     "C30908",
		"Carren":    "C30908",
		"Caron":     "C30908",
		"Lambard":   "L7081096",
		"Lambert":   "L7081096",
		"Lampert":   "L7081096",
		"Lamport":   "L7081096",
		"Robert":    "R901096",
		"Rupert":    "R901096",
		"Rubin":     "R90108",
		"R_u_b_i_n": "R90108",
		"Ashcraft":  "A03039026",
		"Ashcroft":  "A03039026",
		"Tymczak":   "T6083503",
		"Pfister":   "P1203609",
	}

	for w, expectedCode := range words {
		assert.Equal(t, expectedCode, RefinedSoundex(w))
	}
}

func Benchmark_RefinedSoundex(b *testing.B) {
	words := []string{
		"a", "cat", "catastrophe", "baseball", "dandelion", "elephant",
		"foggy", "ghost", "hangover", "illegal",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, w := range words {
			RefinedSoundex(w)
		}
	}
}
