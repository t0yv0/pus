package main

import (
	"github.com/t0yv0/complang/value"
)

type Object struct {
	syms    []value.Symbol
	values  []value.Value
	shownAs string
}

func (o *Object) ShownAs(printed string) *Object {
	return &Object{
		syms:    o.syms,
		values:  o.values,
		shownAs: printed,
	}
}

func (o *Object) With(field string, v value.Value) *Object {
	return &Object{
		syms:    append(o.syms, value.NewSymbol(field)),
		values:  append(o.values, v),
		shownAs: o.shownAs,
	}
}

func (o *Object) Value() value.Value {
	m := map[value.Symbol]value.Value{}
	for i, s := range o.syms {
		v := o.values[i]
		m[s] = v
	}
	return &objectValue{
		Value:   &value.MapValue{Value: m},
		shownAs: o.shownAs,
	}
}

func NewObject() *Object {
	return &Object{}
}

type objectValue struct {
	value.Value
	shownAs string
}

func (ov *objectValue) Run() value.Value {
	return ov
}

func (ov *objectValue) Show() string {
	if ov.shownAs != "" {
		return ov.shownAs
	}
	return "<object>"
}

var _ value.ValueLike = (*objectValue)(nil)
