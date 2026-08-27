package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectCases(t *testing.T) {
	root := t.TempDir()
	cases := reproductionCases(root)

	all, err := selectCases(cases, "")
	if err != nil {
		t.Fatalf("select all cases: %v", err)
	}
	if len(all) != 4 {
		t.Fatalf("selected %d cases, want 4", len(all))
	}

	selected, err := selectCases(cases, "otel-demo-async")
	if err != nil {
		t.Fatalf("select named case: %v", err)
	}
	if len(selected) != 1 || selected[0].Name != "otel-demo-async" {
		t.Fatalf("selected %#v, want otel-demo-async", selected)
	}
}

func TestSelectCasesRejectsUnknownName(t *testing.T) {
	_, err := selectCases(reproductionCases(t.TempDir()), "unknown")
	if err == nil || !strings.Contains(err.Error(), "available:") {
		t.Fatalf("select unknown case error = %v, want available case list", err)
	}
}

func TestRepositoryRootFromEnvironment(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "toolchain-lock.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("REPRODUCTION_ROOT", root)

	got, err := repositoryRoot()
	if err != nil {
		t.Fatalf("repositoryRoot: %v", err)
	}
	want, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("repositoryRoot = %q, want %q", got, want)
	}
}
