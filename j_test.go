package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffJ(t *testing.T) {
	x := mustJ(map[string]interface{}{
		"a": 1,
		"b": "foo",
		"c": []interface{}{1, 2, 3, 4},
	})
	y := mustJ(map[string]interface{}{
		"a": 2,
		"b": "foo",
		"c": []interface{}{0, 1, 2, 5},
	})

	dxy, xneqy := x.diff(y)

	assert.True(t, xneqy)
	assert.Equal(t, map[string]interface{}{
		"a": map[string]interface{}{
			"add": float64(2),
			"rm":  float64(1),
		},
		"c": []interface{}{
			map[string]interface{}{"add": float64(0)},
			float64(1),
			float64(2),
			map[string]interface{}{"rm": float64(3)},
			map[string]interface{}{"rm": float64(4)},
			map[string]interface{}{"add": float64(5)},
		},
	}, dxy.v)

	x = mustJ([]interface{}{
		"id",
	})
	_, xneqx := x.diff(x)
	assert.False(t, xneqx)
}
