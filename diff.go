package main

import (
	"bytes"
	"fmt"
	"strings"
)

func diff[T any](eq func(T, T) bool, xs, ys []T) []edit[T] {
	d := &differ[T]{eq, xs, ys}
	return d.diff()
}

type change int

const (
	insert change = iota
	remove
	keep
)

type edit[T any] struct {
	change  change
	element T
}

type differ[T any] struct {
	eq func(T, T) bool
	xs []T
	ys []T
}

func (d *differ[T]) difflen() *matrix[int] {
	difflen := newMatrix[int](len(d.xs)+1, len(d.ys)+1)
	for xp := len(d.xs); xp >= 0; xp-- {
		for yp := len(d.ys); yp >= 0; yp-- {
			l, _ := d.choose(difflen, xp, yp)
			difflen.set(xp, yp, l)
		}
	}
	return difflen
}

func (d *differ[T]) choose(difflen *matrix[int], xp, yp int) (int, change) {
	xrem := len(d.xs) - xp
	yrem := len(d.ys) - yp
	switch {
	case xrem == 0:
		return yrem, insert
	case yrem == 0:
		return xrem, remove
	}
	l := 1 + difflen.get(xp+1, yp)
	c := remove
	if n := 1 + difflen.get(xp, yp+1); n < l {
		l = n
		c = insert
	}
	if d.eq(d.xs[xp], d.ys[yp]) {
		if n := difflen.get(xp+1, yp+1); n < l {
			l = n
			c = keep
		}
	}
	return l, c
}

func (d *differ[T]) diff() []edit[T] {
	var edits []edit[T]
	difflen, xs, ys := d.difflen(), d.xs, d.ys
	for {
		if len(xs) == 0 {
			for _, y := range ys {
				edits = append(edits, d.insert(y))
			}
			return edits
		}
		if len(ys) == 0 {
			for _, x := range xs {
				edits = append(edits, d.remove(x))
			}
			return edits
		}
		xp, yp := len(d.xs)-len(xs), len(d.ys)-len(ys)
		_, diff := d.choose(difflen, xp, yp)
		switch diff {
		case remove:
			edits, xs = append(edits, d.remove(xs[0])), xs[1:]
		case insert:
			edits, ys = append(edits, d.insert(ys[0])), ys[1:]
		default: // keep
			edits, xs, ys = append(edits, d.keep(xs[0])), xs[1:], ys[1:]
		}
	}
}

func (d *differ[T]) insert(x T) edit[T] {
	return edit[T]{insert, x}
}

func (d *differ[T]) remove(x T) edit[T] {
	return edit[T]{remove, x}
}

func (d *differ[T]) keep(x T) edit[T] {
	return edit[T]{keep, x}
}

type matrix[T any] struct {
	m    int
	n    int
	data []T
}

func newMatrix[T any](m, n int) *matrix[T] {
	return &matrix[T]{m, n, make([]T, m*n)}
}

func (m *matrix[T]) get(i, j int) T {
	return m.data[m.index(i, j)]
}

func (m *matrix[T]) set(i, j int, v T) {
	m.data[m.index(i, j)] = v
}

func (m *matrix[T]) index(i, j int) int {
	return m.n*i + j
}

func diffLines(text1, text2 string) string {
	diffs := diff(func(a, b string) bool { return a == b },
		strings.Split(text1, "\n"),
		strings.Split(text2, "\n"))
	var buf bytes.Buffer
	changed := false
	for i, ed := range diffs {
		switch ed.change {
		case insert:
			fmt.Fprintf(&buf, "+ %s\n", ed.element)
			changed = true
		case remove:
			fmt.Fprintf(&buf, "- %s\n", ed.element)
			changed = true
		case keep:
			nearChanges := false
			d := 5
			for j := i - d; j <= i+d && j >= 0 && j < len(diffs); j++ {
				switch diffs[j].change {
				case insert, remove:
					nearChanges = true
				}
			}
			if nearChanges {
				fmt.Fprintf(&buf, "  %s\n", ed.element)
			}
		}
	}
	if !changed {
		return ""
	}
	return buf.String()
}
