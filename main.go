package main

import (
	"context"
	"encoding/json"
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
	completeSnippet := flag.String("complete", "", "optionally pass some code to complete on; pass '-' for stdin")
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

func readStdinAsLines() ([]string, error) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, l := range strings.Split(string(bytes), "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}
	return lines, nil
}

func initEnv() complang.MutableEnv {
	env := initialEnv()
	menv := complang.NewMutableEnv()
	for k, v := range env {
		menv.Bind(k, v)
	}
	return menv
}

func execCode() error {
	ctx := context.Background()
	lines, err := readStdinAsLines()
	if err != nil {
		return err
	}
	return execLines(ctx, initEnv(), lines)
}

func execLines(ctx context.Context, menv complang.MutableEnv, lines []string) error {
	for _, s := range lines {
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
	menv := initEnv()
	if completeSnippet == "-" {
		lines, err := readStdinAsLines()
		if err != nil {
			return err
		}
		if len(lines) == 0 {
			return fmt.Errorf("Expected more than 1 line")
		}
		completeSnippet = lines[len(lines)-1]
		if err := execLines(ctx, menv, lines[0:len(lines)-1]); err != nil {
			return err
		}
	}
	query, err := parser.ParseQuery(completeSnippet)
	if err != nil {
		return fmt.Errorf("%q: %w", completeSnippet, err)
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
	prefixedMatches := []string{}
	for _, m := range matches {
		prefixedMatches = append(prefixedMatches, fmt.Sprintf("%s%s", pfx, m))
	}
	j, err := json.MarshalIndent(prefixedMatches, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", j)
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
