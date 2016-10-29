package gosim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLDAModel(t *testing.T) {
	ldaModel := NewLDAModel()
	fmt.Println("TestNewLDAModel says 'yo.'")
	assert.NotNil(t, ldaModel)
}
