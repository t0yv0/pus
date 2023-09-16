package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/t0yv0/complang/value"
)

type jview struct {
	Root      j
	Path      jpath
	Current   j
	Transform func(*jview) j
	Extend    func(*jview, *Object) *Object
}

type jpath []jfrag

type jfrag any // Union[int,string]

func (jv *jview) ToValue() value.Value {
	return LazyValue(func() value.Value {
		c := jv.Current
		o := NewObject()
		switch {
		case jv.IsMap():
			m := jv.AsMap()
			for k, v := range m {
				o = o.With(k, v.ToValue())
			}
			if isCompact(c) {
				o = o.ShownAs(c.String())
			} else {
				o = o.RunAs(func() value.Value { return showInFX(c) })
			}
		case jv.IsList():
			ll := jv.AsList()
			for i, v := range ll {
				o = o.With(fmt.Sprintf("_%d", i), v.ToValue())
			}
			if isCompact(c) {
				o = o.ShownAs(c.String())
			} else {
				o = o.RunAs(func() value.Value { return showInFX(c) })
			}
		case jv.IsText():
			o = o.ShownAs(jv.Text())
		default:
			o = o.ShownAs(c.String())
		}
		o = jv.Extend(jv, o)
		return o.Value()
	})
}

func (jv *jview) IsMap() bool {
	_, ok := jv.Current.v.(map[string]interface{})
	return ok
}

func (jv *jview) AsMap() map[string]*jview {
	l, ok := jv.Current.v.(map[string]interface{})
	if !ok {
		return nil
	}
	ret := map[string]*jview{}
	for k, e := range l {
		path := append(jv.Path, jfrag(k))
		ret[k] = &jview{
			Root:      jv.Root,
			Path:      path,
			Current:   j{e},
			Transform: jv.Transform,
			Extend:    jv.Extend,
		}
		ret[k].Current = jv.Transform(ret[k])
	}
	return ret
}

func (jv *jview) IsList() bool {
	_, ok := jv.Current.v.([]interface{})
	return ok
}

func (jv *jview) AsList() []*jview {
	l, ok := jv.Current.v.([]interface{})
	if !ok {
		return nil
	}
	ret := []*jview{}
	for i, e := range l {
		path := append(jv.Path, jfrag(i))
		ret = append(ret, &jview{
			Root:      jv.Root,
			Path:      path,
			Current:   j{e},
			Transform: jv.Transform,
			Extend:    jv.Extend,
		})
	}
	return ret
}

func (jv *jview) IsText() bool {
	_, ok := jv.Current.v.(string)
	return ok
}

func (jv *jview) Text() string {
	s, _ := jv.Current.v.(string)
	return s
}

func isCompact(x j) bool {
	text := x.String()
	lines := strings.Split(text, "\n")
	if len(lines) > 10 {
		return false
	}
	for _, l := range lines {
		if len(l) > 80 {
			return false
		}
	}
	return true
}

func showInFX(x j) value.Value {
	fxPath, err := exec.LookPath("fx")
	if err != nil {
		return &value.ErrorValue{ErrorMessage: err.Error()}
	}
	cmd := exec.Command(fxPath)
	cmd.Stdin = bytes.NewBuffer([]byte(x.String()))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return &value.ErrorValue{ErrorMessage: err.Error()}
	}
	return &value.NullValue{}
}
