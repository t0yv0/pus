package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/t0yv0/complang"
	"github.com/t0yv0/complang/expr"
	"github.com/t0yv0/complang/parser"
	"github.com/t0yv0/complang/repl"
)

func main() {
	completeSnippet := flag.String("complete", "", "optionally pass some code to complete on")
	executeFlag := flag.Bool("execute", false, "optionally pass some code to execute on stdin")
	flag.Parse()

	switch {
	case completeSnippet != nil && *completeSnippet != "":
		err := complete(*completeSnippet)
		if err != nil {
			log.Fatal(err)
		}
	case executeFlag != nil && *executeFlag:
		err := execCode()
		if err != nil {
			log.Fatal(err)
		}
	default:
		startREPL()
	}
}

func execCode() error {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	executeSnippet := string(bytes)
	ctx := context.Background()

	env := initialEnv()
	menv := complang.NewMutableEnv()
	for k, v := range env {
		menv.Bind(k, v)
	}

	for _, s := range strings.Split(executeSnippet, "\n") {
		stmt, err := parser.ParseStmt(s)
		if err != nil {
			return err
		}
		if stmt == nil {
			continue
		}
		expr.EvalStmt(ctx, menv, stmt)
	}
	return nil
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
