package main

import (
	"github.com/t0yv0/complang/value"
)

func strValue(s string) value.Value {
	return &printable{StringValue: value.StringValue{Value: s}}
}

type printable struct {
	value.StringValue
}

func (p *printable) Run() value.Value {
	return p
}

func (p *printable) Show() string {
	return p.StringValue.Value
}
