package main

import (
	"fmt"
	"strings"
)

type env map[string]value

type value interface {
	Message(arg value) value
	Complete(query string) []string
	Show() string
}

type deferredValue struct {
	get func() value
}

func (dv deferredValue) Message(arg value) value {
	return dv.get().Message(arg)
}

func (dv deferredValue) Complete(query string) []string {
	return dv.get().Complete(query)
}

func (dv deferredValue) Show() string {
	return dv.get().Show()
}

type strValue string

func (sv strValue) Message(arg value) value {
	return arg
}

func (sv strValue) Complete(query string) []string {
	return nil
}

func (sv strValue) Show() string {
	return string(sv)
}

func errValue(format string, arg ...any) value {
	return strValue(fmt.Sprintf(format, arg...))
}

type mapValue map[string]value

func (mv mapValue) Message(arg value) value {
	sv, ok := arg.(strValue)
	if !ok {
		return errValue("Error: map only responds to str messages, given %s", arg.Show())
	}
	rv, ok := mv[string(sv)]
	if !ok {
		return errValue("Error: unknown key %s", string(sv))
	}
	return rv
}

func (mv mapValue) Show() string {
	return fmt.Sprintf("map[n=%d]", len(mv))
}

func (mv mapValue) Complete(query string) (matches []string) {
	for k := range mv {
		if !strings.HasPrefix(k, query) {
			continue
		}
		matches = append(matches, k)
	}
	return
}

func newStdEnv() env {
	v := make(env)
	v["schema"] = deferredValue{schemaPrim}
	return v
}

func readEval(env env, rawExpr string, readonly bool) value {
	tokens := strings.Split(rawExpr, " ")

	var assignTo string
	if len(tokens) >= 3 && tokens[1] == "=" {
		assignTo = tokens[0]
		tokens = tokens[2:]
	}

	v, ok := env[tokens[0]]
	if !ok {
		v = strValue(tokens[0])
	}

	for _, m := range tokens[1:] {
		if vv, ok := env[m]; ok {
			v = v.Message(vv)
		} else {
			v = v.Message(strValue(m))
		}
	}

	if assignTo != "" && !readonly {
		env[assignTo] = v
	}
	return v
}
