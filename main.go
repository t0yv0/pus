package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/t0yv0/complang"
	"github.com/t0yv0/complang/expr"
	"github.com/t0yv0/complang/parser"
	"github.com/t0yv0/complang/repl"
)

func main() {
	completeSnippet := flag.String("complete", "", "optionally pass some code to complete on")
	flag.Parse()

	if completeSnippet != nil && *completeSnippet != "" {
		err := complete(*completeSnippet)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		startREPL()
	}
}

func complete(completeSnippet string) error {
	ctx := context.Background()
	query, err := parser.ParseQuery(completeSnippet)
	if err != nil {
		return fmt.Errorf("%q: %w", completeSnippet, err)
	}
	env := initialEnv()
	menv := complang.NewMutableEnv()
	for k, v := range env {
		menv.Bind(k, v)
	}
	matches := []string{}
	qt := query.QueryText()
	expr.EvalQuery(ctx, menv, query, func(_, match string) bool {
		if !strings.HasPrefix(match, qt) {
			return true
		}
		matches = append(matches, match)
		return true
	})
	sort.Strings(matches)
	pfx := completeSnippet[0:query.Offset()]
	for _, m := range matches {
		fmt.Printf("%s%s\n", pfx, m)
	}
	return nil
}

func startREPL() {
	ctx := context.Background()

	err := repl.ReadEvalPrintLoop(ctx, repl.ReadEvalPrintLoopOptions{
		MaxCompletions:     16,
		HistoryFile:        "/tmp/pus.readline.history",
		InitialEnvironment: initialEnv(),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func initialEnv() map[string]complang.Value {
	schema := complang.LazyValue(func() complang.Value {
		s, err := LoadSchema()
		if err != nil {
			return complang.BindValue(err)
		}
		return complang.BindValue(s)
	})
	return map[string]complang.Value{
		"$schema": schema,
	}
}
