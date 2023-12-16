package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/t0yv0/complang"
)

type Schema struct {
	packageSpec *schema.PackageSpec
}

func (s *Schema) Name() string {
	return s.packageSpec.Name
}

func (s *Schema) DisplayName() string {
	return s.packageSpec.DisplayName
}

func (s *Schema) Version() string {
	return s.packageSpec.Version
}

func (s *Schema) Description() string {
	return s.packageSpec.Description
}

func (s *Schema) Keywords() []string {
	return s.packageSpec.Keywords
}

func (s *Schema) Homepage() string {
	return s.packageSpec.Homepage
}

func (s *Schema) License() string {
	return s.packageSpec.License
}

func (s *Schema) Attribution() string {
	return s.packageSpec.Attribution
}

func (s *Schema) Repository() string {
	return s.packageSpec.Repository
}

func (s Schema) LogoURL() string {
	return s.packageSpec.LogoURL
}

func (s *Schema) PluginDownloadURL() string {
	return s.packageSpec.PluginDownloadURL
}

func (s *Schema) Publisher() string {
	return s.packageSpec.Publisher
}

func (s *Schema) AllowedPackageNames() []string {
	return s.packageSpec.AllowedPackageNames
}

func (s *Schema) Language() complang.Value {
	return complang.BindValue(mustJ(s.packageSpec.Language).v)
}

func (s *Schema) Config() complang.Value {
	return complang.BindValue(inlineRefs(s.packageSpec, mustJ(s.packageSpec.Config)).v)
}

func (s *Schema) Meta() complang.Value {
	return complang.BindValue(inlineRefs(s.packageSpec, mustJ(s.packageSpec.Meta)).v)
}

func (s *Schema) Provider() complang.Value {
	return complang.BindValue(inlineRefs(s.packageSpec, mustJ(s.packageSpec.Provider)).v)
}

func (s *Schema) Resources() map[string]complang.Value {
	m := map[string]complang.Value{}
	for k, v := range s.packageSpec.Resources {
		m[k] = newResource(s, k, v)
	}
	return m
}

func (s *Schema) Functions() map[string]complang.Value {
	m := map[string]complang.Value{}
	for k, v := range s.packageSpec.Functions {
		m[k] = newFunction(s, k, v)
	}
	return m
}

func (s *Schema) Types() map[string]complang.Value {
	m := map[string]complang.Value{}
	for k, v := range s.packageSpec.Types {
		m[k] = newType(s, k, v)
	}
	return m
}

type Resource struct {
	schema *Schema
	token  string
	res    schema.ResourceSpec
	desc   string
}

func newResource(s *Schema, tok string, res schema.ResourceSpec) complang.Value {
	copy := res
	copy.Description = ""
	schema := inlineRefs(s.packageSpec, mustJ(copy)).v
	return complang.OverloadedValue(
		complang.BindValue(Resource{s, tok, res, res.Description}),
		complang.BindValue(schema),
	)
}

func (r Resource) Description() string {
	return r.desc
}

func (r Resource) Token() string {
	return r.token
}

type Function struct {
	schema *Schema
	token  string
	fu     schema.FunctionSpec
	desc   string
}

func newFunction(s *Schema, tok string, fu schema.FunctionSpec) complang.Value {
	copy := fu
	copy.Description = ""
	schema := inlineRefs(s.packageSpec, mustJ(copy)).v

	return complang.OverloadedValue(
		complang.BindValue(Function{s, tok, fu, fu.Description}),
		complang.BindValue(schema),
	)
}

func (f Function) Token() string {
	return f.token
}

func (f Function) Description() string {
	return f.desc
}

type Type struct {
	schema *Schema
	token  string
	ty     schema.ComplexTypeSpec
}

func newType(s *Schema, tok string, ty schema.ComplexTypeSpec) complang.Value {
	return complang.OverloadedValue(
		complang.BindValue(Type{s, tok, ty}),
		complang.BindValue(inlineRefs(s.packageSpec, mustJ(ty)).v),
	)
}

func (f Type) Token() string {
	return f.token
}

func LoadSchema() (*Schema, error) {
	spec, err := autoloadSchema()
	if err != nil {
		return nil, err
	}
	return &Schema{spec}, nil
}

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

// Diff between baseline value retrieved as follows, and the given value.
//
//	git show HEAD:..schema.json
func (s *Schema) Diff() complang.Value {
	baseline, err := loadSchemaAtGitRef("HEAD")
	if err != nil {
		return &complang.Error{
			ErrorMessage: fmt.Sprintf("%v", err),
		}
	}
	diff, hasDiff := mustJ(baseline).diff(mustJ(s.packageSpec))
	if hasDiff {
		return complang.BindValue(diff.v)
	}
	return complang.BindValue("No schema changes between HEAD and current version")
}

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
func inlineRefs(s *schema.PackageSpec, v j) j {
	return v.transform(func(inj j) j {
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
