package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

type envCompleter struct {
	env env
}

func (vc *envCompleter) Do(line []rune, pos int) ([][]rune, int) {
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
	v := readEval(vc.env, strings.TrimSpace(beforeLast), true /*readonly*/)
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
	return fuzzyMatch(query, completions, 25)
}

func repl() (finalError error) {
	env := newStdEnv()
	cfg := &readline.Config{
		Prompt:            "\033[31m»\033[0m ",
		HistoryFile:       "/tmp/pus-readline.tmp",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      &envCompleter{env: env},
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
		fmt.Println(readEval(env, line, false /*readonly*/).Show())
	}

	return nil
}
