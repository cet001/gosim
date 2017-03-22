package strdist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevenshtein_Dist(t *testing.T) {
	dist := NewLevenshtein().Dist

	assert.Equal(t, 0, dist("", ""))
	assert.Equal(t, 0, dist("foo", "foo"))

	assert.Equal(t, 3, dist("foo", ""))
	assert.Equal(t, 3, dist("", "foo"))
	assert.Equal(t, 3, dist("foo", "bar"))

	assert.Equal(t, 3, dist("kitten", "sitting"))
	assert.Equal(t, 3, dist("saturday", "sunday"))
}
