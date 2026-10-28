package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHistoricalReferenceReleaseIsSeparateFromCurrentToolchain(t *testing.T) {
	var lock toolchainLock
	lock.Toolchain.Version = "1.3.3"
	if got := referenceRelease(lock); got != "1.3.3" {
		t.Fatalf("legacy lock fallback = %q", got)
	}
	lock.Toolchain.Version = "1.3.4"
	lock.Toolchain.ReferenceRelease = "1.3.3"
	if got := referenceRelease(lock); got != "1.3.3" {
		t.Fatalf("historical reference was silently advanced: %q", got)
	}
}

func TestEvaluationComparisonIgnoresOnlyVerifiedReleaseProvenance(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	historical := write("historical.json", `{"toolchain_release":"1.3.3","availability":0.75,"decision":"pass"}`)
	generated := write("generated.json", `{"toolchain_release":"1.3.4","availability":0.75,"decision":"pass"}`)
	if err := compareEvaluationReference(historical, generated, "1.3.4"); err != nil {
		t.Fatalf("only release provenance differed: %v", err)
	}
	if err := compareEvaluationReference(historical, generated, "1.3.5"); err == nil || !strings.Contains(err.Error(), "provenance mismatch") {
		t.Fatalf("wrong current release was accepted: %v", err)
	}
	generated = write("changed-number.json", `{"toolchain_release":"1.3.4","availability":0.74,"decision":"pass"}`)
	if err := compareEvaluationReference(historical, generated, "1.3.4"); err == nil || !strings.Contains(err.Error(), "reference mismatch") {
		t.Fatalf("changed numerical result was accepted: %v", err)
	}
	generated = write("changed-decision.json", `{"toolchain_release":"1.3.4","availability":0.75,"decision":"fail"}`)
	if err := compareEvaluationReference(historical, generated, "1.3.4"); err == nil || !strings.Contains(err.Error(), "reference mismatch") {
		t.Fatalf("changed decision was accepted: %v", err)
	}
	historical = write("historical-large.json", `{"toolchain_release":"1.3.3","count":9007199254740992}`)
	generated = write("changed-large.json", `{"toolchain_release":"1.3.4","count":9007199254740993}`)
	if err := compareEvaluationReference(historical, generated, "1.3.4"); err == nil || !strings.Contains(err.Error(), "reference mismatch") {
		t.Fatalf("changed large integer was rounded away: %v", err)
	}
}
