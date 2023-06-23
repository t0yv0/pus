package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	repl()
}

type valueCompleter struct {
	env   env
	value value
}

func (vc *valueCompleter) Do(line []rune, pos int) ([][]rune, int) {
	if pos != len(line) {
		return nil, len(line)
	}
	l := runesToStr(line)
	var beforeLast, lastToken string
	i := strings.LastIndexAny(l, " ")
	if i == -1 {
		beforeLast = ""
		lastToken = l
	} else {
		beforeLast = l[0 : i+1]
		lastToken = l[i+1:]
	}
	v := readEval(vc.env, vc.value, beforeLast, true /*readonly*/)
	var out [][]rune
	for _, x := range completeInEnv(vc.env, v, lastToken) {
		if strings.HasPrefix(x, lastToken) {
			out = append(out, strToRunes(strings.TrimPrefix(x, lastToken)))
		}
	}
	return out, len(lastToken)
}

func completeInEnv(env env, v value, query string) []string {
	completions := v.Complete(query)
	for k := range env {
		if strings.HasPrefix(k, query) {
			completions = append(completions, k)
		}
	}
	return completions
}

func repl() (finalError error) {
	env := make(env)
	v := newStdLib()
	cfg := &readline.Config{
		Prompt:            "\033[31m»\033[0m ",
		HistoryFile:       "/tmp/pus-readline.tmp",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      &valueCompleter{value: v, env: env},
		FuncFilterInputRune: func(r rune) (rune, bool) {
			switch r {
			// block CtrlZ feature
			case readline.CharCtrlZ:
				return r, false
			}
			return r, true
		},
	}
	l, err := readline.NewEx(cfg)
	if err != nil {
		return err
	}
	defer func() {
		closeError := l.Close()
		if closeError != nil && finalError != nil {
			finalError = closeError
		}
	}()
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
		fmt.Println(readEval(env, v, line, false /*readonly*/).Show())
	}

	return nil
}
