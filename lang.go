package main

import (
	"strings"
)

type env map[string]value

type value interface {
	Message(arg value) value
	Complete(query string) []string
	Show() string
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

func newStdLib() value {
	return &stdlib{}
}

func readEval(env env, v value, rawExpr string, readonly bool) value {
	tokens := strings.Split(rawExpr, " ")

	var assignTo string
	if len(tokens) >= 3 && tokens[1] == "=" {
		assignTo = tokens[0]
		tokens = tokens[2:]
	}

	for _, m := range tokens {
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

type stdlib struct {
	path []value
}

func (self *stdlib) Message(arg value) value {
	return &stdlib{path: append(self.path, arg)}
}

func (self *stdlib) Complete(prefix string) []string {
	return []string{prefix + "1", prefix + "2", prefix + "3"}
}

func (self *stdlib) Show() string {
	var parts []string
	for _, p := range self.path {
		parts = append(parts, p.Show())
	}
	return strings.Join(parts, ".")
}
