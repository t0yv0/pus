package main

import (
	"log"

	"github.com/t0yv0/complang/repl"
	"github.com/t0yv0/complang/value"
)

func main() {
	err := repl.ReadEvalPrintLoop(repl.ReadEvalPrintLoopOptions{
		HistoryFile: "/tmp/pus.readline.history",
		InitialEnvironment: map[value.Symbol]value.Value{
			value.NewSymbol("$schema"): schemaPrimValue,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
