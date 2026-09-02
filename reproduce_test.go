package main

import (
	"errors"
	"os"
	"os/exec"
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

func TestRunCommandWithExit(t *testing.T) {
	if testing.Short() {
		t.Skip("starts a subprocess")
	}
	name := "go"
	args := []string{"version"}
	if _, err := exec.LookPath(name); err != nil {
		t.Skip("go command is unavailable")
	}
	code, output, err := runCommandWithExit(t.TempDir(), name, args...)
	if err != nil || code != 0 || !strings.HasPrefix(output, "go version") {
		t.Fatalf("runCommandWithExit success = code %d output %q err %v", code, output, err)
	}
}

func TestRunCommandWithExitReportsExitCode(t *testing.T) {
	if testing.Short() {
		t.Skip("starts a subprocess")
	}
	name := "go"
	if _, err := exec.LookPath(name); err != nil {
		t.Skip("go command is unavailable")
	}
	code, _, err := runCommandWithExit(t.TempDir(), name, "tool", "definitely-not-a-go-tool")
	if err != nil || code == 0 {
		t.Fatalf("runCommandWithExit failure = code %d err %v", code, err)
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		t.Fatalf("expected process exit to be represented by code, got %v", err)
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
