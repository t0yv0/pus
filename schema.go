package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

func detectSchemaFile() (string, error) {
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

func schemaValue(spec *schema.PackageSpec) value {
	v := make(mapValue)

	resources := make(mapValue)
	v["rs"] = resources

	for rname := range spec.Resources {
		rname := rname
		resources[rname] = deferredValue{func() value {
			return strValue(pretty(spec.Resources[rname]))
		}}
	}

	fns := make(mapValue)
	v["fn"] = fns

	for fname := range spec.Functions {
		fname := fname
		fns[fname] = deferredValue{func() value {
			return strValue(pretty(spec.Functions[fname]))
		}}
	}

	tys := make(mapValue)
	v["ty"] = tys
	for tname := range spec.Types {
		tname := tname
		tys[tname] = deferredValue{func() value {
			return strValue(pretty(spec.Types[tname]))
		}}
	}

	return v
}

func schemaPrim() value {
	spec, err := autoloadSchema()
	if err != nil {
		return errValue("%v", err)
	}
	return schemaValue(spec)
}

func pretty(x any) string {
	bs, _ := json.MarshalIndent(x, "", "  ")
	return string(bs)
}
