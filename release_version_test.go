package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestSingleReleaseVersion(t *testing.T) {
	read := func(path string) string {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	var lock struct {
		Artifact  map[string]any `json:"artifact"`
		Toolchain struct {
			Version string `json:"version"`
		} `json:"toolchain"`
	}
	if err := json.Unmarshal([]byte(read("toolchain-lock.json")), &lock); err != nil {
		t.Fatal(err)
	}
	if _, found := lock.Artifact["version"]; found {
		t.Fatal("artifact must not have an independent version")
	}
	version := lock.Toolchain.Version
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(version) {
		t.Fatalf("invalid release version %q", version)
	}
	want := "ghcr.io/mb3r-lab/mb3r-software-complex-repro:v" + version
	imagePattern := regexp.MustCompile(`ghcr\.io/mb3r-lab/mb3r-software-complex-repro:[A-Za-z0-9_.-]+`)
	for _, path := range []string{"compose.yaml", "README.md", "docs/manual.md", "docs/troubleshooting.md"} {
		matches := imagePattern.FindAllString(read(path), -1)
		if len(matches) == 0 {
			t.Errorf("%s has no versioned image instruction", path)
		}
		for _, actual := range matches {
			if actual != want {
				t.Errorf("%s: image %s, want %s", path, actual, want)
			}
		}
	}
	for _, path := range []string{"Dockerfile", ".github/workflows/reproduce.yml"} {
		content := read(path)
		if strings.Contains(content, "ARTIFACT_VERSION") || strings.Contains(content, ".artifact.version") || strings.Contains(content, "artifact_version") {
			t.Errorf("%s still has an independent artifact version", path)
		}
	}
	for _, path := range []string{"reference/manifest.json", "reference/evaluation/decision.json", "reference/evaluation/validation-matrix.json", "reference/operation/result.json"} {
		var reference struct {
			Version string `json:"toolchain_release"`
		}
		if err := json.Unmarshal([]byte(read(path)), &reference); err != nil {
			t.Fatal(err)
		}
		if reference.Version != version {
			t.Errorf("%s: release %s, want %s", path, reference.Version, version)
		}
	}
	read("release/v" + version + ".md")
}
