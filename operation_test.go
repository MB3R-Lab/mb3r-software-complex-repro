package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestOperationStoredEvidence(t *testing.T) {
	stamp, err := validateOperationEvidence(filepath.Join("cases", "operation-semantics", "evidence"))
	if err != nil {
		t.Fatal(err)
	}
	if stamp < "2026-10-15T00:00:00Z" {
		t.Fatalf("wrong collection timestamp %s", stamp)
	}
}

func TestOperationResponseCriterion(t *testing.T) {
	for _, tc := range []struct {
		status float64
		body   string
		want   bool
	}{
		{200, `[{"id":"p1"}]`, true}, {200, `[]`, false}, {200, `null`, false},
		{200, `[{"id":""}]`, false}, {200, `[{"id":"p1"},{}]`, false},
		{200, `invalid`, false}, {500, `[{"id":"p1"}]`, false},
	} {
		if got := operationResponseSuccess(tc.status, tc.body); got != tc.want {
			t.Errorf("%v %s: got %v want %v", tc.status, tc.body, got, tc.want)
		}
	}
}

func TestOperationWilson(t *testing.T) {
	ci := operationWilson(.95, 40000)
	if ci["lower_bound"] >= .95 || ci["upper_bound"] <= .95 {
		t.Fatal(ci)
	}
}

func TestOperationReferenceIsReleaseSpecific(t *testing.T) {
	root := t.TempDir()
	if got, want := operationReferencePath(root, "1.3.3"), filepath.Join(root, "reference", "operation", "result.json"); got != want {
		t.Fatalf("historical reference path = %q, want %q", got, want)
	}
	if got, want := operationReferencePath(root, "1.3.4"), filepath.Join(root, "reference", "operation-v1.3.4", "result.json"); got != want {
		t.Fatalf("new release reference path = %q, want %q", got, want)
	}
	if got := operationReferencePath(root, "1.3.5"); got == operationReferencePath(root, "1.3.3") {
		t.Fatal("new release silently fell back to the historical operation reference")
	}
}

func TestCurrentOperationReferenceProvenance(t *testing.T) {
	lock, err := loadToolchainLock("toolchain-lock.json")
	if err != nil {
		t.Fatal(err)
	}
	if lock.Toolchain.Version != "1.3.4" {
		t.Skip("this provenance record applies to the 1.3.4 candidate")
	}
	var provenance struct {
		ToolchainRelease    string            `json:"toolchain_release"`
		ReferenceKind       string            `json:"reference_kind"`
		SourceRunID         string            `json:"source_run_id"`
		SourceHeadCommit    string            `json:"source_head_commit"`
		SourceArtifact      string            `json:"source_artifact"`
		SourceComponents    map[string]string `json:"source_components"`
		ResultSHA256        string            `json:"result_sha256"`
		HistoricalSHA256    string            `json:"historical_operation_reference_sha256_on_linux_runner"`
	}
	provenanceRaw, err := os.ReadFile("reference/operation-v1.3.4/provenance.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(provenanceRaw, &provenance); err != nil {
		t.Fatal(err)
	}
	if provenance.ToolchainRelease != lock.Toolchain.Version || provenance.ReferenceKind != "full_graph_operation_replay" || provenance.SourceRunID == "" || provenance.SourceHeadCommit == "" || provenance.SourceArtifact == "" {
		t.Fatalf("incomplete current operation provenance: %+v", provenance)
	}
	for name, component := range lock.Components {
		if provenance.SourceComponents[name] != component.Source.Commit {
			t.Fatalf("%s source differs from frozen operation provenance", name)
		}
	}
	for path, want := range map[string]string{
		"reference/operation-v1.3.4/result.json": provenance.ResultSHA256,
		"reference/operation/result.json": provenance.HistoricalSHA256,
	} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n")))
		if got := fmt.Sprintf("sha256:%x", sum); got != want {
			t.Fatalf("%s: normalized file hash %s, want %s", path, got, want)
		}
	}
}
