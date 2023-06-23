package main

import (
	"fmt"
	"unicode/utf8"
)

func strToRunes(s string) []rune {
	if s == "" {
		return nil
	}
	var runes []rune
	bytes := []byte(s)
	for {
		r, n := utf8.DecodeRune(bytes)
		if r == utf8.RuneError {
			panic(fmt.Sprintf("utf8 decoding failed: %q rem=%q", s, string(bytes)))
		}
		runes = append(runes, r)
		bytes = bytes[n:]
		if len(bytes) == 0 {
			return runes
		}
	}
}

func runesToStr(runes []rune) string {
	if len(runes) == 0 {
		return ""
	}
	var bytes []byte
	for _, r := range runes {
		bytes = utf8.AppendRune(bytes, r)
	}
	return string(bytes)
}
