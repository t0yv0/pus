package main

import (
	"context"
	"log"

	"github.com/t0yv0/complang"
	"github.com/t0yv0/complang/repl"
)

func main() {
	ctx := context.Background()
	schema := complang.LazyValue(func() complang.Value {
		s, err := LoadSchema()
		if err != nil {
			return complang.BindValue(err)
		}
		return complang.BindValue(s)
	})
	err := repl.ReadEvalPrintLoop(ctx, repl.ReadEvalPrintLoopOptions{
		MaxCompletions: 16,
		HistoryFile:    "/tmp/pus.readline.history",
		InitialEnvironment: map[string]complang.Value{
			"$schema": schema,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
