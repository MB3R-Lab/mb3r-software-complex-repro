package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const discoveredAt = "2026-08-21T00:00:00Z"

var referenceFiles = []string{
	"procrustes/report.json",
	"procrustes/collection-plan.json",
	"procrustes/bering.overlay.yaml",
	"procrustes/summary.md",
	"bering/model.json",
	"bering/snapshot.json",
	"bering/model-quality.json",
	"bering/snapshot-quality.json",
	"sheaft/report.json",
	"sheaft/model.json",
	"sheaft/summary.md",
	"analytics.json",
	"summary.md",
}

type componentPin struct {
	Version string `json:"version"`
	Source  struct {
		Commit string `json:"commit"`
	} `json:"source"`
}

type toolchainLock struct {
	Toolchain struct {
		Version string `json:"version"`
	} `json:"toolchain"`
	Components map[string]componentPin `json:"components"`
}

type caseSpec struct {
	Name        string
	SourceAlias string
	Dir         string
}

type referenceManifest struct {
	SchemaVersion    string                       `json:"schema_version"`
	ToolchainRelease string                       `json:"toolchain_release"`
	DiscoveredAt     string                       `json:"discovered_at"`
	Normalization    []string                     `json:"normalization"`
	Cases            map[string]map[string]string `json:"cases"`
}

func main() {
	update := flag.Bool("update-reference", false, "replace checked-in semantic reference outputs")
	flag.Parse()
	if err := run(*update); err != nil {
		fmt.Fprintln(os.Stderr, "reproduce-paper:", err)
		os.Exit(1)
	}
}

func run(update bool) error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	manifest, err := loadToolchainLock(filepath.Join(root, "toolchain-lock.json"))
	if err != nil {
		return err
	}
	repos := map[string]string{
		"procrustes": repositoryPath(root, "PROCRUSTES_REPO", "Procrustes"),
		"bering":     repositoryPath(root, "BERING_REPO", "Bering"),
		"sheaft":     repositoryPath(root, "SHEAFT_REPO", "Sheaft"),
	}
	for _, name := range []string{"procrustes", "bering", "sheaft"} {
		pin, ok := manifest.Components[name]
		if !ok {
			return fmt.Errorf("toolchain manifest is missing %s", name)
		}
		actual, err := commandOutput(repos[name], "git", "rev-parse", "HEAD")
		if err != nil {
			return fmt.Errorf("read %s revision: %w", name, err)
		}
		if actual != pin.Source.Commit {
			return fmt.Errorf("%s revision mismatch: got %s, want %s", name, actual, pin.Source.Commit)
		}
		fmt.Printf("pin-ok %-11s version=%s commit=%s\n", name, pin.Version, actual[:12])
	}

	tmp := filepath.Join(root, ".tmp", "paper-reproduction")
	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return err
	}
	bins := map[string]string{}
	for _, item := range []struct{ name, command string }{
		{"procrustes", "./cmd/procrustes"},
		{"bering", "./cmd/bering"},
		{"sheaft", "./cmd/sheaft"},
	} {
		binary := filepath.Join(binDir, item.name+executableSuffix())
		if err := runCommand("build "+item.name, repos[item.name], "go", "build", "-trimpath", "-o", binary, item.command); err != nil {
			return err
		}
		bins[item.name] = binary
	}

	cases := []caseSpec{
		{Name: "otel-demo", SourceAlias: "paper", Dir: filepath.Join(root, "cases", "otel-demo")},
		{Name: "social-network", SourceAlias: "compose", Dir: filepath.Join(root, "cases", "social-network")},
	}
	allHashes := map[string]map[string]string{}
	workRoot := filepath.Join(tmp, "work")
	for _, c := range cases {
		caseWork := filepath.Join(workRoot, c.Name)
		if err := runCase(c, caseWork, bins); err != nil {
			return err
		}
		first, err := semanticHashes(caseWork)
		if err != nil {
			return err
		}
		if err := runCase(c, caseWork, bins); err != nil {
			return err
		}
		second, err := semanticHashes(caseWork)
		if err != nil {
			return err
		}
		if err := compareHashes(c.Name+" repeated run", first, second); err != nil {
			return err
		}
		allHashes[c.Name] = second
		fmt.Printf("stable %-14s files=%d\n", c.Name, len(second))
	}

	referenceRoot := filepath.Join(root, "reference")
	wanted := referenceManifest{
		SchemaVersion:    "1.0.0",
		ToolchainRelease: manifest.Toolchain.Version,
		DiscoveredAt:     discoveredAt,
		Normalization: []string{
			"remove generated_at and recompute_duration_ms",
			"remove machine-local artifact.path and input_artifact.path",
			"remove the non-semantic Procrustes instrumentation evidence field selector (source document is retained)",
			"canonicalize JSON object ordering and whitespace",
		},
		Cases: allHashes,
	}
	if update {
		for _, c := range cases {
			if err := writeReferenceCase(filepath.Join(workRoot, c.Name), filepath.Join(referenceRoot, c.Name)); err != nil {
				return err
			}
		}
		if err := writeJSON(filepath.Join(referenceRoot, "manifest.json"), wanted); err != nil {
			return err
		}
		fmt.Println("reference outputs updated")
		return nil
	}

	var current referenceManifest
	if err := readJSON(filepath.Join(referenceRoot, "manifest.json"), &current); err != nil {
		return fmt.Errorf("read reference manifest (run make reproduce-paper-update once): %w", err)
	}
	if current.ToolchainRelease != wanted.ToolchainRelease || current.DiscoveredAt != wanted.DiscoveredAt {
		return errors.New("reference manifest belongs to a different toolchain release")
	}
	for _, c := range cases {
		if err := compareHashes(c.Name+" reference", current.Cases[c.Name], wanted.Cases[c.Name]); err != nil {
			return err
		}
	}
	fmt.Printf("reproduce-paper-ok toolchain=%s cases=%d\n", manifest.Toolchain.Version, len(cases))
	return nil
}

func runCase(c caseSpec, work string, bins map[string]string) error {
	if err := os.RemoveAll(work); err != nil {
		return err
	}
	procrustesOut := filepath.Join(work, "procrustes")
	beringOut := filepath.Join(work, "bering")
	sheaftOut := filepath.Join(work, "sheaft")
	for _, dir := range []string{procrustesOut, beringOut, sheaftOut} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if err := runCommand(c.Name+" procrustes", c.Dir, bins["procrustes"],
		"analyze", "--source", c.SourceAlias+"="+filepath.Join(c.Dir, "static"),
		"--goal", filepath.Join(c.Dir, "goal.yaml"), "--out-dir", procrustesOut, "--fail-on-blocked"); err != nil {
		return err
	}
	modelPath := filepath.Join(beringOut, "model.json")
	if err := runCommand(c.Name+" bering", c.Dir, bins["bering"],
		"discover", "--input", filepath.Join(c.Dir, "topology-api.yaml"),
		"--overlay", filepath.Join(procrustesOut, "bering.overlay.yaml"),
		"--discovered-at", discoveredAt,
		"--out", modelPath,
		"--snapshot-out", filepath.Join(beringOut, "snapshot.json"),
		"--quality-out", filepath.Join(beringOut, "model-quality.json"),
		"--snapshot-quality-out", filepath.Join(beringOut, "snapshot-quality.json")); err != nil {
		return err
	}
	if err := runCommand(c.Name+" sheaft", c.Dir, bins["sheaft"],
		"run", "--model", modelPath, "--analysis", filepath.Join(c.Dir, "analysis.yaml"), "--out-dir", sheaftOut); err != nil {
		return err
	}
	return writeAnalytics(c, work)
}

func writeAnalytics(c caseSpec, work string) error {
	var preflight, model, report map[string]any
	if err := readJSON(filepath.Join(work, "procrustes", "report.json"), &preflight); err != nil {
		return err
	}
	if err := readJSON(filepath.Join(work, "bering", "model.json"), &model); err != nil {
		return err
	}
	if err := readJSON(filepath.Join(work, "sheaft", "report.json"), &report); err != nil {
		return err
	}
	services := arrayAt(model, "services")
	edges := arrayAt(model, "edges")
	endpoints := arrayAt(model, "endpoints")
	asyncEdges := 0
	for _, raw := range edges {
		if edge, ok := raw.(map[string]any); ok && stringAt(edge, "kind") == "async" {
			asyncEdges++
		}
	}
	analytics := map[string]any{
		"schema_version": "1.0.0",
		"case":           c.Name,
		"procrustes": map[string]any{
			"goal_status":              stringAt(preflight, "goal_status"),
			"services_assessed":        len(arrayAt(preflight, "services")),
			"blockers":                 len(arrayAt(preflight, "blockers")),
			"remediation_items":        len(arrayAt(preflight, "remediation")),
			"capability_status_counts": statusCounts(arrayAt(preflight, "capabilities")),
		},
		"bering": map[string]any{
			"services": len(services), "edges": len(edges), "async_edges": asyncEdges, "endpoints": len(endpoints),
			"confidence": nestedAt(model, "metadata", "confidence"),
		},
		"sheaft": map[string]any{
			"summary":          report["summary"],
			"endpoint_results": report["endpoint_results"],
			"sweeps":           report["sweeps"],
		},
	}
	if c.Name == "social-network" {
		var legacy map[string]any
		if err := readJSON(filepath.Join(c.Dir, "legacy-reference.json"), &legacy); err != nil {
			return err
		}
		deltas := legacyDeltas(report, legacy)
		if accepted, _ := deltas["within_tolerance"].(bool); !accepted {
			return fmt.Errorf("social-network result exceeds the archived comparison tolerance: max absolute delta %.6f", asFloat(deltas["max_absolute_delta"]))
		}
		comparison := map[string]any{"reference": legacy, "deltas": deltas}
		analytics["legacy_contextual_comparison"] = comparison
	}
	if err := writeJSON(filepath.Join(work, "analytics.json"), analytics); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(work, "summary.md"), []byte(markdownSummary(c.Name, analytics)), 0o644)
}

func legacyDeltas(report, legacy map[string]any) map[string]any {
	deltas := map[string]any{}
	summary, _ := report["summary"].(map[string]any)
	aggregateDelta := numberAt(summary, "weighted_overall_availability") - numberAt(legacy, "weighted_aggregate")
	deltas["weighted_aggregate"] = aggregateDelta
	maxAbsolute := math.Abs(aggregateDelta)
	legacyEndpoints, _ := legacy["endpoints"].(map[string]any)
	endpointDeltas := map[string]any{}
	for _, raw := range arrayAt(report, "endpoint_results") {
		endpoint, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id := stringAt(endpoint, "endpoint_id")
		if previous, ok := legacyEndpoints[id]; ok {
			delta := numberAt(endpoint, "availability") - asFloat(previous)
			endpointDeltas[id] = delta
			maxAbsolute = math.Max(maxAbsolute, math.Abs(delta))
		}
	}
	deltas["endpoints"] = endpointDeltas
	deltas["tolerance"] = 0.015
	deltas["max_absolute_delta"] = maxAbsolute
	deltas["within_tolerance"] = maxAbsolute <= 0.015
	return deltas
}

func markdownSummary(name string, analytics map[string]any) string {
	p, _ := analytics["procrustes"].(map[string]any)
	b, _ := analytics["bering"].(map[string]any)
	s, _ := analytics["sheaft"].(map[string]any)
	summary, _ := s["summary"].(map[string]any)
	var out strings.Builder
	fmt.Fprintf(&out, "# %s reference summary\n\n", name)
	fmt.Fprintf(&out, "- Procrustes goal status: `%s`; assessed services: %.0f; blockers: %.0f.\n", stringAt(p, "goal_status"), asFloat(p["services_assessed"]), asFloat(p["blockers"]))
	fmt.Fprintf(&out, "- Bering model: %.0f services, %.0f edges (%.0f async), %.0f endpoints.\n", asFloat(b["services"]), asFloat(b["edges"]), asFloat(b["async_edges"]), asFloat(b["endpoints"]))
	fmt.Fprintf(&out, "- Sheaft p=0.30 weighted availability: %.6f.\n\n", numberAt(summary, "weighted_overall_availability"))
	out.WriteString("| Endpoint | Availability |\n|---|---:|\n")
	for _, raw := range arrayAt(s, "endpoint_results") {
		endpoint, ok := raw.(map[string]any)
		if ok {
			fmt.Fprintf(&out, "| %s | %.6f |\n", stringAt(endpoint, "endpoint_id"), numberAt(endpoint, "availability"))
		}
	}
	return out.String()
}

func semanticHashes(root string) (map[string]string, error) {
	result := map[string]string{}
	for _, rel := range referenceFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		raw, err := normalizedFile(path)
		if err != nil {
			return nil, fmt.Errorf("normalize %s: %w", rel, err)
		}
		sum := sha256.Sum256(raw)
		result[rel] = "sha256:" + hex.EncodeToString(sum[:])
	}
	return result, nil
}

func normalizedFile(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if filepath.Ext(path) != ".json" {
		return raw, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	normalize(value, "")
	return json.MarshalIndent(value, "", "  ")
}

func normalize(value any, parent string) {
	switch typed := value.(type) {
	case map[string]any:
		delete(typed, "generated_at")
		delete(typed, "recompute_duration_ms")
		if parent == "artifact" || parent == "input_artifact" {
			delete(typed, "path")
		}
		if id, ok := typed["id"].(string); ok && strings.HasSuffix(id, ".instrumentation") {
			if evidence, ok := typed["evidence"].([]any); ok {
				for _, raw := range evidence {
					if item, ok := raw.(map[string]any); ok {
						delete(item, "field")
					}
				}
			}
		}
		for key, child := range typed {
			normalize(child, key)
		}
	case []any:
		for _, child := range typed {
			normalize(child, parent)
		}
	}
}

func writeReferenceCase(work, destination string) error {
	if err := os.RemoveAll(destination); err != nil {
		return err
	}
	for _, rel := range referenceFiles {
		raw, err := normalizedFile(filepath.Join(work, filepath.FromSlash(rel)))
		if err != nil {
			return err
		}
		target := filepath.Join(destination, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if filepath.Ext(target) == ".json" {
			raw = append(raw, '\n')
		}
		if err := os.WriteFile(target, raw, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func compareHashes(label string, want, got map[string]string) error {
	keys := make([]string, 0, len(want)+len(got))
	seen := map[string]bool{}
	for key := range want {
		seen[key] = true
		keys = append(keys, key)
	}
	for key := range got {
		if !seen[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	var differences []string
	for _, key := range keys {
		if want[key] != got[key] {
			differences = append(differences, fmt.Sprintf("%s: want=%s got=%s", key, want[key], got[key]))
		}
	}
	if len(differences) > 0 {
		return fmt.Errorf("%s mismatch:\n  %s", label, strings.Join(differences, "\n  "))
	}
	return nil
}

func repositoryRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("cannot locate reproduction source")
	}
	return filepath.Dir(file), nil
}

func repositoryPath(root, envName, sibling string) string {
	if value := strings.TrimSpace(os.Getenv(envName)); value != "" {
		return value
	}
	for _, base := range []string{filepath.Dir(root), filepath.Dir(filepath.Dir(root))} {
		candidate := filepath.Join(base, sibling)
		if info, err := os.Stat(filepath.Join(candidate, ".git")); err == nil && info.IsDir() {
			return candidate
		}
	}
	return filepath.Join(filepath.Dir(root), sibling)
}

func executableSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func runCommand(label, dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w\n%s", label, err, strings.TrimSpace(output.String()))
	}
	return nil
}

func commandOutput(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	raw, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(raw)))
	}
	return strings.TrimSpace(string(raw)), nil
}

func loadToolchainLock(path string) (toolchainLock, error) {
	var value toolchainLock
	if err := readJSON(path, &value); err != nil {
		return value, err
	}
	return value, nil
}

func readJSON(path string, destination any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	return decoder.Decode(destination)
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func statusCounts(values []any) map[string]int {
	counts := map[string]int{}
	for _, raw := range values {
		if value, ok := raw.(map[string]any); ok {
			counts[stringAt(value, "status")]++
		}
	}
	return counts
}

func arrayAt(value map[string]any, key string) []any { out, _ := value[key].([]any); return out }
func nestedAt(value map[string]any, first, second string) any {
	child, _ := value[first].(map[string]any)
	return child[second]
}
func stringAt(value map[string]any, key string) string  { out, _ := value[key].(string); return out }
func numberAt(value map[string]any, key string) float64 { return asFloat(value[key]) }
func asFloat(value any) float64 {
	switch typed := value.(type) {
	case json.Number:
		out, _ := typed.Float64()
		return out
	case float64:
		return typed
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	}
	return 0
}
