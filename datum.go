package main

import (
	"encoding/json"
	"fmt"

	"github.com/t0yv0/complang/value"
	"gopkg.in/yaml.v3"
)

// Wraps a JSON-serializable value into an object for user interaction.
func datum(raw any) value.Value {
	v, err := jsonify(raw)
	if err != nil {
		return fromError(err)
	}

	o := NewObject()

	switch v := v.(type) {
	case string:
		o = o.ShownAs(v)
	default:
		s, err := showAny(v)
		if err != nil {
			return fromError(err)
		}
		o = o.ShownAs(s)
	}

	switch v := v.(type) {
	case map[string]interface{}:
		for key, subv := range v {
			o = o.With(key, datum(subv))
		}
	case []interface{}:
		for i, subv := range v {
			o = o.With(fmt.Sprintf("_%d", i), datum(subv))
		}
	}

	return o.Value()
}

func jsonify(v any) (any, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func showAny(v any) (string, error) {
	y, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(y), nil
}

func fromError(err error) value.Value {
	return &value.ErrorValue{ErrorMessage: fmt.Sprintf("%v", err)}
}
