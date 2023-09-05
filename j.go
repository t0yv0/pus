package main

import (
	"reflect"

	"encoding/json"
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

func (x j) diff(y j) (j, bool) {
	switch xv := x.v.(type) {
	case []interface{}:
		switch yv := y.v.(type) {
		case []interface{}:
			d := []interface{}{}
			diffs := diff(reflect.DeepEqual, xv, yv)
			xeqy := true
			for _, ed := range diffs {
				switch ed.change {
				case insert:
					d = append(d, j{ed.element}.added().v)
					xeqy = false
				case remove:
					d = append(d, j{ed.element}.removed().v)
					xeqy = false
				case keep:
					d = append(d, j{ed.element}.v)
				}
			}
			if xeqy {
				return j{}, false
			}
			return mustJ(d), true
		default:
			return x.changed(y), true
		}
	case map[string]interface{}:
		switch yv := y.v.(type) {
		case map[string]interface{}:
			diffmap := map[string]interface{}{}
			eq := true
			for k, xvv := range xv {
				yvv, ok := yv[k]
				if !ok {
					eq = false
					diffmap[k] = j{xvv}.removed().v
				} else {
					xdy, xneqy := j{xvv}.diff(j{yvv})
					if xneqy {
						eq = false
						diffmap[k] = xdy.v
					}
				}
			}
			for k, yvv := range yv {
				if _, ok := xv[k]; !ok {
					eq = false
					diffmap[k] = j{yvv}.added().v
				}
			}
			if eq {
				return j{}, false
			}
			return mustJ(diffmap), true
		default:
			return x.changed(y), true
		}
	default:
		if reflect.DeepEqual(x.v, y.v) {
			return j{}, false
		}
		return x.changed(y), true
	}
}

func (x j) removed() j {
	return mustJ(map[string]interface{}{"-": x.v})
}

func (x j) added() j {
	return mustJ(map[string]interface{}{"+": x.v})
}

func (x j) changed(y j) j {
	return mustJ(map[string]interface{}{"-": x.v, "+": y.v})
}
