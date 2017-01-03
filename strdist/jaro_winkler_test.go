package strdist

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleJaroWinklerDist() {
	dist := NewJaroWinklerDist()
	fmt.Println(dist.CalcString("andrew", "andrew"))
	fmt.Println(dist.CalcString("martha", "marhta"))
	fmt.Println(dist.CalcString("jones", "johnson"))
	fmt.Println(dist.CalcString("foo", "bar"))
	// Output:
	// 1
	// 0.9611111111111111
	// 0.8323809523809523
	// 0
}

func TestJaroWinklerDist_CalcString(t *testing.T) {
	dist := NewJaroWinklerDist()

	assert.Equal(t, 0.0, dist.CalcString("", ""))
	assert.Equal(t, 0.0, dist.CalcString("abc", ""))
	assert.Equal(t, 0.0, dist.CalcString("", "abc"))
	assert.Equal(t, 0.0, dist.CalcString("abc", "xyz"))

	assert.Equal(t, 1.0, dist.CalcString("abc", "abc"))

	// See examples 1 and 2 from https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance
	assert.Equal(t, 0.9611111111111111, dist.CalcString("martha", "marhta"))
	assert.Equal(t, 0.8133333333333332, dist.CalcString("dixon", "dicksonx"))

	// See examples from http://alias-i.com/lingpipe/docs/api/com/aliasi/spell/JaroWinklerDistance.html
	assert.Equal(t, 0.8323809523809523, dist.CalcString("jones", "johnson"))
}

func TestJaroWinklerDist_Calc(t *testing.T) {
	dist := NewJaroWinklerDist()

	assert.Equal(t, 0.0, dist.Calc([]byte(""), []byte("")))
	assert.Equal(t, 0.0, dist.Calc([]byte("abc"), []byte("")))
	assert.Equal(t, 0.0, dist.Calc([]byte(""), []byte("abc")))
	assert.Equal(t, 0.0, dist.Calc([]byte("abc"), []byte("xyz")))

	assert.Equal(t, 1.0, dist.Calc([]byte("abc"), []byte("abc")))

	// See examples 1 and 2 from https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance
	assert.Equal(t, 0.9611111111111111, dist.Calc([]byte("martha"), []byte("marhta")))
	assert.Equal(t, 0.8133333333333332, dist.Calc([]byte("dixon"), []byte("dicksonx")))

	// See examples from http://alias-i.com/lingpipe/docs/api/com/aliasi/spell/JaroWinklerDistance.html
	assert.Equal(t, 0.8323809523809523, dist.Calc([]byte("jones"), []byte("johnson")))
}

func Benchmark_JaroWinklerDist_CalcString(b *testing.B) {
	dist := NewJaroWinklerDist()

	s1values := []string{"martha", "dixon", "apple", "constitution", "mississippi"}
	s2values := []string{"marhta", "dicksonx", "microsoft", "intervention", "misanthrope"}

	numValues := len(s1values)
	i := 0

	calcDist := func() {
		dist.CalcString(s1values[i], s2values[i])

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

func Benchmark_JaroWinklerDist_Calc(b *testing.B) {
	dist := NewJaroWinklerDist()

	s1values := [][]byte{
		[]byte("martha"),
		[]byte("dixon"),
		[]byte("apple"),
		[]byte("constitution"),
		[]byte("mississippi"),
	}

	s2values := [][]byte{
		[]byte("marhta"),
		[]byte("dicksonx"),
		[]byte("microsoft"),
		[]byte("intervention"),
		[]byte("misanthrope"),
	}

	numValues := len(s1values)
	i := 0

	calcDist := func() {
		dist.Calc(s1values[i], s2values[i])

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
