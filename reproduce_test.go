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

func TestCheckedInCriticalityReports(t *testing.T) {
	root, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range reproductionCases(root) {
		t.Run(c.Name, func(t *testing.T) {
			var report map[string]any
			path := filepath.Join(root, "reference", c.Name, "criticality", "report.json")
			if err := readJSON(path, &report); err != nil {
				t.Fatalf("read criticality reference: %v", err)
			}
			analysis, err := criticalityAnalysis(c, report)
			if err != nil {
				t.Fatalf("validate criticality reference: %v", err)
			}
			wantEndpoints := 4
			if c.Family == "social-network" {
				wantEndpoints = 3
			}
			if got := len(arrayAt(analysis, "endpoints")); got != wantEndpoints {
				t.Fatalf("criticality endpoint count = %d, want %d", got, wantEndpoints)
			}
		})
	}
}

func TestSelectCriticalityWalkthroughCase(t *testing.T) {
	cases := reproductionCases(t.TempDir())
	got, err := selectCriticalityWalkthroughCase(cases, "")
	if err != nil {
		t.Fatalf("select default criticality walkthrough: %v", err)
	}
	if got.Name != "otel-demo-async" {
		t.Fatalf("default criticality walkthrough = %q, want otel-demo-async", got.Name)
	}

	got, err = selectCriticalityWalkthroughCase(cases, "social-network-repl")
	if err != nil {
		t.Fatalf("select named criticality walkthrough: %v", err)
	}
	if got.Name != "social-network-repl" {
		t.Fatalf("named criticality walkthrough = %q, want social-network-repl", got.Name)
	}
}

func TestNumericOracleCheck(t *testing.T) {
	pass := newNumericCheck("pass", "example", 0.9, 0.895, 0.01)
	if pass.Status != "PASS" {
		t.Fatalf("numeric check status = %s, want PASS", pass.Status)
	}
	fail := newNumericCheck("fail", "example", 0.9, 0.85, 0.01)
	if fail.Status != "FAIL" {
		t.Fatalf("numeric check status = %s, want FAIL", fail.Status)
	}
}

func TestSetDifferenceReturnsExplicitEmptySlice(t *testing.T) {
	got := setDifference([]string{"a", "b"}, []string{"a", "b"})
	if got == nil || len(got) != 0 {
		t.Fatalf("setDifference = %#v, want non-nil empty slice", got)
	}
}

func TestCheckedInEvaluationReferences(t *testing.T) {
	var decision map[string]any
	if err := readJSON(filepath.Join("reference", "evaluation", "decision.json"), &decision); err != nil {
		t.Fatal(err)
	}
	stages := arrayAt(decision, "stages")
	if len(stages) != 3 {
		t.Fatalf("decision stages = %d, want 3", len(stages))
	}
	wantDecisions := []string{"pass", "fail", "pass"}
	for index, raw := range stages {
		stage, _ := raw.(map[string]any)
		if got := stringAt(stage, "decision"); got != wantDecisions[index] {
			t.Fatalf("decision stage %d = %q, want %q", index, got, wantDecisions[index])
		}
	}

	var validation map[string]any
	if err := readJSON(filepath.Join("reference", "evaluation", "validation-matrix.json"), &validation); err != nil {
		t.Fatal(err)
	}
	summary, _ := validation["summary"].(map[string]any)
	if got := int(numberAt(summary, "passed")); got != 15 {
		t.Fatalf("validation passed = %d, want 15", got)
	}
}
