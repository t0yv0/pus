package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

const maxMatchItems = 32

var packageSpec *schema.PackageSpec

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

func fuzzyMatch(query string, collection []string) (out []string) {
	matches := fuzzy.RankFindFold(query, collection)
	sort.Sort(matches)
	for _, m := range matches {
		out = append(out, m.Target)
	}
	return out
}

func findResources(query string) []string {
	query = strings.TrimPrefix(query, "resource ")
	if packageSpec == nil {
		return nil
	}
	var tokens []string
	for k := range packageSpec.Resources {
		tokens = append(tokens, k)
	}
	matches := fuzzyMatch(query, tokens)
	if len(matches) > maxMatchItems {
		matches = matches[0:maxMatchItems]
	}
	return matches
}

func findFunctions(query string) []string {
	query = strings.TrimPrefix(query, "function ")
	if packageSpec == nil {
		return nil
	}
	var tokens []string
	for k := range packageSpec.Functions {
		tokens = append(tokens, k)
	}
	matches := fuzzyMatch(query, tokens)
	if len(matches) > maxMatchItems {
		matches = matches[0:maxMatchItems]
	}
	return matches
}

func findTypes(query string) []string {
	query = strings.TrimPrefix(query, "type ")
	if packageSpec == nil {
		return nil
	}
	var tokens []string
	for k := range packageSpec.Types {
		tokens = append(tokens, k)
	}
	matches := fuzzyMatch(query, tokens)
	if len(matches) > maxMatchItems {
		matches = matches[0:maxMatchItems]
	}
	return matches
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("help"),
	readline.PcItem("resource", readline.PcItemDynamic(findResources)),
	readline.PcItem("function", readline.PcItemDynamic(findFunctions)),
	readline.PcItem("type", readline.PcItemDynamic(findTypes)),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func detectSchemaFile() (string, bool) {
	cmd := exec.Command("git", "ls-files", "**/schema.json")
	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err == nil {
		s := out.String()
		s = strings.TrimSpace(s)
		if s != "" {
			return s, true
		}
	}
	return "", false
}

func autoloadSchema() error {
	if f, ok := detectSchemaFile(); ok {
		bytes, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(bytes, &packageSpec); err != nil {
			return err
		}
	}
	return nil
}

func main1() {
	if err := autoloadSchema(); err != nil {
		// ignore
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/pus-readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == "help":
			usage(l.Stderr())
		case strings.HasPrefix(line, "resource"):
			r := strings.TrimPrefix(line, "resource ")
			fmt.Printf("resource %q:\n%s\n", r, pretty(packageSpec.Resources[r]))
		case strings.HasPrefix(line, "function"):
			f := strings.TrimPrefix(line, "function ")
			fmt.Printf("function %q:\n%s\n", f, pretty(packageSpec.Functions[f]))
		case strings.HasPrefix(line, "type"):
			t := strings.TrimPrefix(line, "type ")
			fmt.Printf("type %q:\n%s\n", t, pretty(packageSpec.Types[t]))
		default:
			log.Println("unrecognized command, try typing 'help':", strconv.Quote(line))
		}
	}
}

func pretty(x any) string {
	bs, _ := json.MarshalIndent(x, "", "  ")
	return string(bs)
}
