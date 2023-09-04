package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	eq := func(a, b byte) bool {
		return a == b
	}
	input := []byte(`mario`)
	dd := diff(eq, input, []byte(`darius`))
	assert.Equal(t, remove, dd[0].change)
	assert.Equal(t, insert, dd[1].change)
	assert.Equal(t, keep, dd[2].change)
	assert.Equal(t, keep, dd[3].change)
	assert.Equal(t, keep, dd[4].change)
	assert.Equal(t, remove, dd[5].change)
	assert.Equal(t, insert, dd[6].change)
	assert.Equal(t, insert, dd[7].change)
}
