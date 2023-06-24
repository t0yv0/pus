package main

import (
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

func fuzzyMatch(query string, collection []string, maxItems int) (out []string) {
	matches := fuzzy.RankFindFold(query, collection)
	sort.Sort(matches)
	for _, m := range matches {
		out = append(out, m.Target)
	}
	if len(out) >= maxItems {
		return out[0:maxItems]
	}
	return out
}
