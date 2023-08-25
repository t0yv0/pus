package main

import (
	"sync"

	"github.com/t0yv0/complang/value"
)

func LazyValue(f func() value.Value) value.Value {
	var once sync.Once
	var actual value.Value
	get := func() value.Value {
		once.Do(func() {
			actual = f()
		})
		return actual
	}
	return DeferredValue(get)
}

func DeferredValue(f func() value.Value) value.Value {
	return &value.CustomValue{ValueLike: &deferredValue{f}}
}

type deferredValue struct {
	v func() value.Value
}

func (dv *deferredValue) Message(arg value.Value) value.Value {
	return dv.v().Message(arg)
}

func (dv *deferredValue) CompleteSymbol(query value.Symbol) []value.Symbol {
	return dv.v().CompleteSymbol(query)
}

func (dv *deferredValue) Run() value.Value {
	return dv.v().Run()
}

func (dv *deferredValue) Show() string {
	return dv.v().Show()
}

var _ value.ValueLike = (*deferredValue)(nil)
