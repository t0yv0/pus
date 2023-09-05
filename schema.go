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
	return detectSchemaGitPath()
}

func detectSchemaGitPath() (string, error) {
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

func loadSchemaAtGitRef(ref string) (*schema.PackageSpec, error) {
	schemaPath, err := detectSchemaGitPath()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", ref, schemaPath))
	cmd.Stdout = &buf
	cmd.Stderr = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var s schema.PackageSpec
	if err := json.Unmarshal(buf.Bytes(), &s); err != nil {
		return nil, err
	}
	return &s, nil
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

func functionValue(name string, pkg *schema.PackageSpec, spec schema.FunctionSpec) value.Value {
	return LazyValue(func() value.Value {
		desc := spec.Description
		spec.Description = ""
		return NewObject().
			With("desc", strValue(desc)).
			With("rawSchema", datum(spec)).
			With("schema", datum(inlineRefs(pkg, spec))).
			ShownAs(fmt.Sprintf("<function:%s>", name)).
			Value()
	})
}

func typeValue(name string, pkg *schema.PackageSpec, spec schema.ComplexTypeSpec) value.Value {
	return LazyValue(func() value.Value {
		desc := spec.Description
		spec.Description = ""
		return NewObject().
			With("desc", strValue(desc)).
			With("rawSchema", datum(spec)).
			With("schema", datum(inlineRefs(pkg, spec))).
			ShownAs(fmt.Sprintf("<type:%s>", name)).
			Value()
	})
}

func resourceValue(name string, pkg *schema.PackageSpec, res schema.ResourceSpec) value.Value {
	return LazyValue(func() value.Value {
		desc := res.Description
		res.Description = ""
		return NewObject().
			With("desc", strValue(desc)).
			With("rawSchema", datum(res)).
			With("schema", datum(inlineRefs(pkg, res))).
			ShownAs(fmt.Sprintf("<resource:%s>", name)).
			Value()
	})
}

func functionsValue(spec *schema.PackageSpec) value.Value {
	o := NewObject()
	for name, fspec := range spec.Functions {
		o = o.With(name, functionValue(name, spec, fspec))
	}
	o = o.ShownAs("<functions>")
	return o.Value()
}

func resourcesValue(spec *schema.PackageSpec) value.Value {
	o := NewObject()
	for name, rspec := range spec.Resources {
		o = o.With(name, resourceValue(name, spec, rspec))
	}
	o = o.ShownAs("<resources>")
	return o.Value()
}

func typesValue(spec *schema.PackageSpec) value.Value {
	o := NewObject()
	for name, tspec := range spec.Types {
		o = o.With(name, typeValue(name, spec, tspec))
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
		With("diff", LazyValue(func() value.Value { return schemaDiff(spec) })).
		ShownAs(showPackageSpec(spec)).
		Value()
}

// Diff between baseline value retrieved as follows, and the given value.
//
//	git show HEAD:..schema.json
func schemaDiff(given *schema.PackageSpec) value.Value {
	baseline, err := loadSchemaAtGitRef("HEAD")
	if err != nil {
		return &value.ErrorValue{
			ErrorMessage: fmt.Sprintf("%v", err),
		}
	}
	diff, hasDiff := mustJ(baseline).diff(mustJ(given))
	if hasDiff {
		return datum(diff.v)
	}
	return datum("No schema changes between HEAD and current version")
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

// Inlines local type references for easier viewing and diffing.
//
// Given input:
//
//	filter:
//	    $ref: '#/types/cloudflare:index/getZonesFilter:getZonesFilter'
//
// Will lookup type definition for "cloudflare:index/getZonesFilter:getZonesFilter" in s.Types
// and inline it:
//
//	filter:
//	    $ref: '#/types/cloudflare:index/getZonesFilter:getZonesFilter'
//	    type: object
//	    properties:
//	        accountId:
//	            description: |
//	                The account identifier to target for the resource.
//	            type: string
func inlineRefs(s *schema.PackageSpec, v any) any {
	return mustJ(v).transform(func(inj j) j {
		switch injv := inj.v.(type) {
		case map[string]interface{}:
			ref, hasRef := injv["$ref"]
			if !hasRef {
				return inj
			}
			refs, ok := ref.(string)
			if !ok {
				return inj
			}
			if !strings.HasPrefix(refs, "#/types/") {
				return inj
			}
			tref := strings.TrimPrefix(refs, "#/types/")
			t, gotT := s.Types[tref]
			if !gotT {
				return inj
			}
			copy := map[string]interface{}{}
			for k, v := range injv {
				copy[k] = v
			}
			for k, v := range mustJ(t).v.(map[string]interface{}) {
				copy[k] = v
			}
			return j{copy}
		default:
			return inj
		}
	})
}
