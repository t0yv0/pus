package main

import (
	"encoding/json"
)

type j struct {
	v interface{}
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

func (inj j) clone() j {
	return mustJ(inj.v)
}

func (inj j) transform(f func(j) j) j {
	switch jv := inj.v.(type) {
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
		return f(inj)
	}
}
