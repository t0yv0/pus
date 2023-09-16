package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
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
			Root:      r,
			Current:   r,
			Path:      nil,
			Transform: transform,
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

func transform(jv *jview) j {
	if len(jv.Path) > 3 {
		var s schema.PackageSpec
		if err := json.Unmarshal([]byte(jv.Root.String()), &s); err == nil {
			re := inlineRefs(&s, jv.Current)
			fmt.Sprintln(re.String())
			return re
		}
	}
	return jv.Current
}
