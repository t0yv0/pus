package main

import (
	"fmt"
	"log"

	"github.com/t0yv0/complang/repl"
	"github.com/t0yv0/complang/value"
)

func main() {
	schemaPrim := LazyValue(func() value.Value {
		s, err := autoloadSchema()
		if err != nil {
			return &value.ErrorValue{
				ErrorMessage: fmt.Sprintf("%v", err),
			}
		}
		r := mustJ(s)
		jv := &jview{
			Root:    r,
			Current: r,
			Path:    nil,
			Transform: func(jv *jview) j {
				return jv.Current
			},
			Extend: func(jv *jview, o *Object) *Object {
				return o
			},
		}
		return jv.ToValue()
	})

	err := repl.ReadEvalPrintLoop(repl.ReadEvalPrintLoopOptions{
		HistoryFile: "/tmp/pus.readline.history",
		InitialEnvironment: map[value.Symbol]value.Value{
			value.NewSymbol("$schema"): schemaPrim,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
