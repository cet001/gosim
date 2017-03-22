package strdist

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleJaroWinklerDist() {
	jaroWinkler := NewJaroWinkler()
	fmt.Println(jaroWinkler.Dist("andrew", "andrew"))
	fmt.Println(jaroWinkler.Dist("martha", "marhta"))
	fmt.Println(jaroWinkler.Dist("jones", "johnson"))
	fmt.Println(jaroWinkler.Dist("foo", "bar"))
	// Output:
	// 1
	// 0.9611111111111111
	// 0.8323809523809523
	// 0
}

func TestJaroWinklerDist_CalcString(t *testing.T) {
	jaroWinkler := NewJaroWinkler()

	assert.Equal(t, 0.0, jaroWinkler.Dist("", ""))
	assert.Equal(t, 0.0, jaroWinkler.Dist("abc", ""))
	assert.Equal(t, 0.0, jaroWinkler.Dist("", "abc"))
	assert.Equal(t, 0.0, jaroWinkler.Dist("abc", "xyz"))

	assert.Equal(t, 1.0, jaroWinkler.Dist("abc", "abc"))

	// See examples 1 and 2 from https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance
	assert.Equal(t, 0.9611111111111111, jaroWinkler.Dist("martha", "marhta"))
	assert.Equal(t, 0.8133333333333332, jaroWinkler.Dist("dixon", "dicksonx"))

	// See examples from http://alias-i.com/lingpipe/docs/api/com/aliasi/spell/JaroWinklerDistance.html
	assert.Equal(t, 0.8323809523809523, jaroWinkler.Dist("jones", "johnson"))
}

func Benchmark_JaroWinkler_Dist(b *testing.B) {
	s1values := []string{"martha", "dixon", "apple", "constitution", "mississippi"}
	s2values := []string{"marhta", "dicksonx", "microsoft", "intervention", "misanthrope"}
	numValues := len(s1values)
	i := 0

	jaroWinkler := NewJaroWinkler()

	calcDist := func() {
		// dist.CalcString(s1values[i], s2values[i])
		jaroWinkler.Dist(s1values[i], s2values[i])
		i++
		if i == numValues {
			i = 0
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calcDist()
	}
}
