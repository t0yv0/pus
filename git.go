package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/coreos/go-semver/semver"
)

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

func detectGitTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "--list")
	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Error calling `git tag --list`: %w", err)
	}
	tags := []string{}
	s := out.String()
	s = strings.TrimSpace(s)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		_, err := semver.NewVersion(strings.TrimPrefix(line, "v"))
		if err != nil {
			continue
		}
		tags = append(tags, line)
	}
	sort.Slice(tags, func(i, j int) bool {
		vi := semver.New(strings.TrimPrefix(tags[i], "v"))
		vj := semver.New(strings.TrimPrefix(tags[j], "v"))
		return !vi.LessThan(*vj)
	})
	return tags, nil
}
