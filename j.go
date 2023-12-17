package main

import (
	"encoding/json"
	"reflect"

	"github.com/t0yv0/godifft"
)

// Newtype-ish wrapper for JSON-like values, that is nulls, strings, numbers, bools and slices or
// stringly-keyed maps over them.
type j struct {
	v interface{}
}

func (x j) String() string {
	b, err := json.MarshalIndent(x.v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mustJ(v any) j {
	j, err := newJ(v)
	if err != nil {
		panic(err)
	}
	return j
}

func newJ(v any) (j, error) {
	if vv, ok := v.(j); ok {
		return vv, nil
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return j{}, err
	}
	var decoded any
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return j{}, err
	}
	return j{decoded}, nil
}

func (x j) clone() j {
	return mustJ(x.v)
}

func (x j) transform(f func(j) j) j {
	switch jv := x.v.(type) {
	case []interface{}:
		tjv := []interface{}{}
		for _, subj := range jv {
			tjv = append(tjv, j{subj}.transform(f).v)
		}
		return f(j{tjv})
	case map[string]interface{}:
		tjv := map[string]interface{}{}
		for k, v := range jv {
			tjv[k] = j{v}.transform(f).v
		}
		return f(j{tjv})
	default:
		return f(x)
	}
}

type jDiffer struct{}

var _ godifft.Differ[any, any] = (*jDiffer)(nil)

func (jd *jDiffer) Added(x any) any   { return mustJ(x).added().v }
func (jd *jDiffer) Removed(x any) any { return mustJ(x).removed().v }

func (jd *jDiffer) Diff(x, y any) (any, bool) {
	if reflect.DeepEqual(x, y) {
		return nil, false
	}
	return mustJ(x).changed(mustJ(y)).v, true
}

func (x j) diff(y j) (j, bool) {
	d, ok := godifft.DiffTree(&jDiffer{}, reflect.DeepEqual, x.v, y.v)
	if !ok {
		return mustJ(nil), false
	}
	return mustJ(d), true
}

func (x j) removed() j {
	return mustJ(map[string]interface{}{"rm": x.v})
}

func (x j) added() j {
	return mustJ(map[string]interface{}{"add": x.v})
}

func (x j) changed(y j) j {
	return mustJ(map[string]interface{}{"rm": x.v, "add": y.v})
}
