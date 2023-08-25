package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/t0yv0/complang/value"
)

func detectSchemaFile() (string, error) {
	if _, err := os.ReadFile("schema.json"); err == nil {
		return "schema.json", nil
	}

	cmd := exec.Command("git", "ls-files", "**schema.json")
	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error calling `git ls-files **schema.json`: %w", err)
	}
	s := out.String()
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("No schema found by calling `git ls-files **schema.json`")
	}
	return s, nil
}

func autoloadSchema() (*schema.PackageSpec, error) {
	var packageSpec schema.PackageSpec
	f, err := detectSchemaFile()
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("Error reading schema file %q: %w", f, err)
	}
	if err := json.Unmarshal(bytes, &packageSpec); err != nil {
		return nil, fmt.Errorf("Error unmarshalling schema: %w", err)
	}
	return &packageSpec, nil
}

func functionValue(spec schema.FunctionSpec) value.Value {
	return strValue(pretty(spec))
}

func typeValue(spec schema.ComplexTypeSpec) value.Value {
	return strValue(pretty(spec))
}

func resourceValue(name string, res schema.ResourceSpec) value.Value {
	desc := res.Description
	res.Description = ""
	return NewObject().
		With("desc", strValue(desc)).
		With("shape", strValue(pretty(res))).
		ShownAs(fmt.Sprintf("<resource:%s>", name)).
		Value()
}

func functionsValue(spec *schema.PackageSpec) value.Value {
	o := NewObject()
	for name, spec := range spec.Functions {
		o = o.With(name, functionValue(spec))
	}
	o = o.ShownAs("<functions>")
	return o.Value()
}

func resourcesValue(spec *schema.PackageSpec) value.Value {
	o := NewObject()
	for name, spec := range spec.Resources {
		o = o.With(name, resourceValue(name, spec))
	}
	o = o.ShownAs("<resources>")
	return o.Value()
}

func typesValue(spec *schema.PackageSpec) value.Value {
	o := NewObject()
	for name, spec := range spec.Types {
		o = o.With(name, typeValue(spec))
	}
	o = o.ShownAs("<types>")
	return o.Value()
}

func showPackageSpec(spec *schema.PackageSpec) string {
	return fmt.Sprintf("<schema:%drs/%dfn/%dty>",
		len(spec.Resources), len(spec.Functions), len(spec.Types))
}

func schemaValue(spec *schema.PackageSpec) value.Value {
	return NewObject().
		With("rs", resourcesValue(spec)).
		With("fn", functionsValue(spec)).
		With("ty", typesValue(spec)).
		ShownAs(showPackageSpec(spec)).
		Value()
}

func pretty(x any) string {
	bs, _ := json.MarshalIndent(x, "", "  ")
	return string(bs)
}

var (
	schemaPrimValue = LazyValue(func() value.Value {
		s, err := autoloadSchema()
		if err != nil {
			return &value.ErrorValue{
				ErrorMessage: fmt.Sprintf("%v", err),
			}
		}
		return schemaValue(s)
	})
)
