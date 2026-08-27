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

var aggregateReferenceFiles = []string{
	"comparison.json",
	"FULL_COMPARISON.md",
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
	Name                string
	Family              string
	Mode                string
	SourceAlias         string
	StaticDir           string
	Dir                 string
	TopologyFile        string
	AnalysisFile        string
	PublishedModelField string
}

type referenceManifest struct {
	SchemaVersion    string                       `json:"schema_version"`
	ToolchainRelease string                       `json:"toolchain_release"`
	DiscoveredAt     string                       `json:"discovered_at"`
	Normalization    []string                     `json:"normalization"`
	Cases            map[string]map[string]string `json:"cases"`
	Aggregate        map[string]string            `json:"aggregate"`
}

type runOptions struct {
	Update   bool
	Prebuilt bool
	CaseName string
	Repeat   int
	WorkRoot string
}

func main() {
	update := flag.Bool("update-reference", false, "replace checked-in semantic reference outputs")
	acceptWorkDir := flag.String("accept-work-dir", "", "normalize completed CI work outputs into the checked-in references without rerunning tools")
	prebuilt := flag.Bool("prebuilt", false, "use Procrustes, Bering, and Sheaft binaries from PATH or *_BIN variables")
	caseName := flag.String("case", "", "run one named case instead of the complete four-case reproduction")
	repeat := flag.Int("repeat", 2, "number of identical executions used to check deterministic output")
	workRoot := flag.String("work-dir", "", "directory for generated outputs (default: .tmp/paper-reproduction/work)")
	listCases := flag.Bool("list-cases", false, "list available case names and exit")
	flag.Parse()
	if *listCases {
		root, err := repositoryRoot()
		if err != nil {
			fmt.Fprintln(os.Stderr, "reproduce-paper:", err)
			os.Exit(1)
		}
		for _, c := range reproductionCases(root) {
			fmt.Println(c.Name)
		}
		return
	}
	if *update && strings.TrimSpace(*acceptWorkDir) != "" {
		fmt.Fprintln(os.Stderr, "reproduce-paper: --update-reference and --accept-work-dir are mutually exclusive")
		os.Exit(1)
	}
	if *repeat < 1 {
		fmt.Fprintln(os.Stderr, "reproduce-paper: --repeat must be at least 1")
		os.Exit(1)
	}
	if *update && strings.TrimSpace(*caseName) != "" {
		fmt.Fprintln(os.Stderr, "reproduce-paper: --update-reference requires the complete case set")
		os.Exit(1)
	}
	if *update && *prebuilt {
		fmt.Fprintln(os.Stderr, "reproduce-paper: --update-reference is disabled for packaged prebuilt binaries")
		os.Exit(1)
	}
	var err error
	if strings.TrimSpace(*acceptWorkDir) != "" {
		err = acceptReferenceWork(*acceptWorkDir)
	} else {
		err = run(runOptions{
			Update:   *update,
			Prebuilt: *prebuilt,
			CaseName: strings.TrimSpace(*caseName),
			Repeat:   *repeat,
			WorkRoot: strings.TrimSpace(*workRoot),
		})
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "reproduce-paper:", err)
		os.Exit(1)
	}
}

func acceptReferenceWork(workRoot string) error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	lock, err := loadToolchainLock(filepath.Join(root, "toolchain-lock.json"))
	if err != nil {
		return err
	}
	workRoot, err = filepath.Abs(workRoot)
	if err != nil {
		return err
	}
	allHashes := map[string]map[string]string{}
	referenceRoot := filepath.Join(root, "reference")
	for _, c := range reproductionCases(root) {
		caseWork := filepath.Join(workRoot, c.Name)
		if err := writeAnalytics(c, caseWork); err != nil {
			return fmt.Errorf("refresh %s analytics: %w", c.Name, err)
		}
		hashes, err := semanticHashes(caseWork)
		if err != nil {
			return fmt.Errorf("accept %s: %w", c.Name, err)
		}
		if err := writeReferenceCase(caseWork, filepath.Join(referenceRoot, c.Name)); err != nil {
			return err
		}
		allHashes[c.Name] = hashes
	}
	if err := writeAggregateComparison(workRoot, reproductionCases(root)); err != nil {
		return err
	}
	aggregateHashes, err := semanticHashesFor(workRoot, aggregateReferenceFiles)
	if err != nil {
		return err
	}
	if err := writeReferenceFiles(workRoot, referenceRoot, aggregateReferenceFiles); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(referenceRoot, "manifest.json"), newReferenceManifest(lock.Toolchain.Version, allHashes, aggregateHashes)); err != nil {
		return err
	}
	fmt.Printf("reference outputs accepted from %s\n", workRoot)
	return nil
}

func run(options runOptions) error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	manifest, err := loadToolchainLock(filepath.Join(root, "toolchain-lock.json"))
	if err != nil {
		return err
	}
	tmp := filepath.Join(root, ".tmp", "paper-reproduction")
	bins, err := prepareBinaries(root, tmp, manifest, options.Prebuilt)
	if err != nil {
		return err
	}

	allCases := reproductionCases(root)
	cases, err := selectCases(allCases, options.CaseName)
	if err != nil {
		return err
	}
	allHashes := map[string]map[string]string{}
	workRoot := filepath.Join(tmp, "work")
	if options.WorkRoot != "" {
		workRoot, err = filepath.Abs(options.WorkRoot)
		if err != nil {
			return fmt.Errorf("resolve work directory: %w", err)
		}
	}
	for _, c := range cases {
		caseWork := filepath.Join(workRoot, c.Name)
		var first, current map[string]string
		for iteration := 1; iteration <= options.Repeat; iteration++ {
			if err := runCase(c, caseWork, bins); err != nil {
				return err
			}
			current, err = semanticHashes(caseWork)
			if err != nil {
				return err
			}
			if iteration == 1 {
				first = current
				continue
			}
			if err := compareHashes(c.Name+" repeated run", first, current); err != nil {
				return err
			}
		}
		allHashes[c.Name] = current
		fmt.Printf("case-ok %-23s files=%d repeats=%d\n", c.Name, len(current), options.Repeat)
	}
	complete := len(cases) == len(allCases)
	aggregateHashes := map[string]string{}
	if complete {
		if err := writeAggregateComparison(workRoot, cases); err != nil {
			return err
		}
		aggregateHashes, err = semanticHashesFor(workRoot, aggregateReferenceFiles)
		if err != nil {
			return err
		}
	}

	referenceRoot := filepath.Join(root, "reference")
	wanted := newReferenceManifest(manifest.Toolchain.Version, allHashes, aggregateHashes)
	if options.Update {
		for _, c := range cases {
			if err := writeReferenceCase(filepath.Join(workRoot, c.Name), filepath.Join(referenceRoot, c.Name)); err != nil {
				return err
			}
		}
		if err := writeReferenceFiles(workRoot, referenceRoot, aggregateReferenceFiles); err != nil {
			return err
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
	if complete {
		if err := compareHashes("aggregate reference", current.Aggregate, wanted.Aggregate); err != nil {
			return err
		}
	}
	fmt.Printf("reproduce-paper-ok toolchain=%s cases=%d repeats=%d\n", manifest.Toolchain.Version, len(cases), options.Repeat)
	return nil
}

func selectCases(cases []caseSpec, name string) ([]caseSpec, error) {
	if name == "" {
		return cases, nil
	}
	for _, c := range cases {
		if c.Name == name {
			return []caseSpec{c}, nil
		}
	}
	available := make([]string, 0, len(cases))
	for _, c := range cases {
		available = append(available, c.Name)
	}
	return nil, fmt.Errorf("unknown case %q (available: %s)", name, strings.Join(available, ", "))
}

func prepareBinaries(root, tmp string, manifest toolchainLock, prebuilt bool) (map[string]string, error) {
	components := []struct {
		name, command, repositoryEnv, sibling, binaryEnv string
	}{
		{"procrustes", "./cmd/procrustes", "PROCRUSTES_REPO", "Procrustes", "PROCRUSTES_BIN"},
		{"bering", "./cmd/bering", "BERING_REPO", "Bering", "BERING_BIN"},
		{"sheaft", "./cmd/sheaft", "SHEAFT_REPO", "Sheaft", "SHEAFT_BIN"},
	}
	bins := map[string]string{}
	if prebuilt {
		for _, item := range components {
			pin, ok := manifest.Components[item.name]
			if !ok {
				return nil, fmt.Errorf("toolchain manifest is missing %s", item.name)
			}
			binary := strings.TrimSpace(os.Getenv(item.binaryEnv))
			if binary == "" {
				var err error
				binary, err = exec.LookPath(item.name + executableSuffix())
				if err != nil {
					return nil, fmt.Errorf("locate prebuilt %s binary: %w", item.name, err)
				}
			}
			if info, err := os.Stat(binary); err != nil || info.IsDir() {
				return nil, fmt.Errorf("prebuilt %s binary is not readable at %s", item.name, binary)
			}
			bins[item.name] = binary
			fmt.Printf("binary-ok %-11s version=%s commit=%s\n", item.name, pin.Version, pin.Source.Commit[:12])
		}
		return bins, nil
	}

	repos := map[string]string{}
	for _, item := range components {
		repos[item.name] = repositoryPath(root, item.repositoryEnv, item.sibling)
		pin, ok := manifest.Components[item.name]
		if !ok {
			return nil, fmt.Errorf("toolchain manifest is missing %s", item.name)
		}
		actual, err := commandOutput(repos[item.name], "git", "rev-parse", "HEAD")
		if err != nil {
			return nil, fmt.Errorf("read %s revision: %w", item.name, err)
		}
		if actual != pin.Source.Commit {
			return nil, fmt.Errorf("%s revision mismatch: got %s, want %s", item.name, actual, pin.Source.Commit)
		}
		fmt.Printf("pin-ok %-11s version=%s commit=%s\n", item.name, pin.Version, actual[:12])
	}
	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return nil, err
	}
	for _, item := range components {
		binary := filepath.Join(binDir, item.name+executableSuffix())
		if err := runCommand("build "+item.name, repos[item.name], "go", "build", "-trimpath", "-o", binary, item.command); err != nil {
			return nil, err
		}
		bins[item.name] = binary
	}
	return bins, nil
}

func reproductionCases(root string) []caseSpec {
	return []caseSpec{
		{Name: "otel-demo-all-blocking", Family: "otel-demo", Mode: "all-blocking", SourceAlias: "paper", StaticDir: "static", Dir: filepath.Join(root, "cases", "otel-demo"), TopologyFile: "topology-all-blocking.yaml", AnalysisFile: "analysis.yaml", PublishedModelField: "model_all_blocking"},
		{Name: "otel-demo-async", Family: "otel-demo", Mode: "async", SourceAlias: "paper", StaticDir: "static", Dir: filepath.Join(root, "cases", "otel-demo"), TopologyFile: "topology-api.yaml", AnalysisFile: "analysis.yaml", PublishedModelField: "model_async"},
		{Name: "social-network-norepl", Family: "social-network", Mode: "norepl", SourceAlias: "compose", StaticDir: "static-norepl", Dir: filepath.Join(root, "cases", "social-network"), TopologyFile: "topology-norepl.yaml", AnalysisFile: "analysis.yaml", PublishedModelField: "model_mean"},
		{Name: "social-network-repl", Family: "social-network", Mode: "repl", SourceAlias: "compose", StaticDir: "static", Dir: filepath.Join(root, "cases", "social-network"), TopologyFile: "topology-api.yaml", AnalysisFile: "analysis.yaml", PublishedModelField: "model_mean"},
	}
}

func newReferenceManifest(toolchainVersion string, cases map[string]map[string]string, aggregate map[string]string) referenceManifest {
	return referenceManifest{
		SchemaVersion:    "1.0.0",
		ToolchainRelease: toolchainVersion,
		DiscoveredAt:     discoveredAt,
		Normalization: []string{
			"remove generated_at and recompute_duration_ms",
			"remove machine-local artifact.path and input_artifact.path",
			"replace machine-local Bering discovery and overlay references with case-stable references",
			"replace identifiers and digests derived from machine-local artifact bytes with valid zero-value identifiers",
			"remove the non-semantic Procrustes instrumentation evidence field selector (source document is retained)",
			"canonicalize JSON object ordering and whitespace",
		},
		Cases:     cases,
		Aggregate: aggregate,
	}
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
		"analyze", "--source", c.SourceAlias+"="+filepath.Join(c.Dir, c.StaticDir),
		"--goal", filepath.Join(c.Dir, "goal.yaml"), "--out-dir", procrustesOut, "--fail-on-blocked"); err != nil {
		return err
	}
	modelPath := filepath.Join(beringOut, "model.json")
	if err := runCommand(c.Name+" bering", c.Dir, bins["bering"],
		"discover", "--input", filepath.Join(c.Dir, c.TopologyFile),
		"--overlay", filepath.Join(procrustesOut, "bering.overlay.yaml"),
		"--discovered-at", discoveredAt,
		"--out", modelPath,
		"--snapshot-out", filepath.Join(beringOut, "snapshot.json"),
		"--quality-out", filepath.Join(beringOut, "model-quality.json"),
		"--snapshot-quality-out", filepath.Join(beringOut, "snapshot-quality.json")); err != nil {
		return err
	}
	if err := runCommand(c.Name+" sheaft", c.Dir, bins["sheaft"],
		"run", "--model", modelPath, "--analysis", filepath.Join(c.Dir, c.AnalysisFile), "--out-dir", sheaftOut); err != nil {
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
		"family":         c.Family,
		"mode":           c.Mode,
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
	var publication map[string]any
	if err := readJSON(filepath.Join(c.Dir, "publication.json"), &publication); err != nil {
		return fmt.Errorf("read %s publication metadata: %w", c.Name, err)
	}
	analytics["publication"] = publication
	comparison, err := fullPublishedComparison(c, report, publication)
	if err != nil {
		return err
	}
	analytics["published_comparison"] = comparison
	if err := writeJSON(filepath.Join(work, "analytics.json"), analytics); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(work, "summary.md"), []byte(markdownSummary(c.Name, analytics)), 0o644)
}

func fullPublishedComparison(c caseSpec, report, publication map[string]any) (map[string]any, error) {
	reported, _ := publication["reported_results"].(map[string]any)
	profiles := profilesByName(report)
	rows := make([]any, 0, 5)
	for _, raw := range arrayAt(reported, "rows") {
		publishedRow, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if c.Family == "social-network" && stringAt(publishedRow, "mode") != c.Mode {
			continue
		}
		failureFraction := numberAt(publishedRow, "failure_fraction")
		profileName := fmt.Sprintf("p%02d", int(math.Round(failureFraction*100)))
		profile, ok := profiles[profileName]
		if !ok {
			return nil, fmt.Errorf("%s: Sheaft report is missing profile %s", c.Name, profileName)
		}
		simulation, _ := profile["simulation"].(map[string]any)
		actual := numberAt(simulation, "weighted_aggregate")
		publishedModel := numberAt(publishedRow, c.PublishedModelField)
		liveField := "live"
		modelSDField := ""
		liveSDField := ""
		if c.Family == "social-network" {
			liveField = "live_mean"
			modelSDField = "model_sd"
			liveSDField = "live_sd"
		}
		publishedLive := numberAt(publishedRow, liveField)
		row := map[string]any{
			"profile":                        profileName,
			"mode":                           c.Mode,
			"failure_fraction":               failureFraction,
			"fixed_k_failures":               numberAt(simulation, "fixed_k_failures"),
			"trials":                         numberAt(simulation, "trials"),
			"sheaft":                         actual,
			"published_model":                publishedModel,
			"published_live":                 publishedLive,
			"delta_from_published_model":     actual - publishedModel,
			"absolute_delta_published_model": math.Abs(actual - publishedModel),
			"delta_from_published_live":      actual - publishedLive,
			"absolute_delta_published_live":  math.Abs(actual - publishedLive),
		}
		if modelSDField != "" {
			row["published_model_sd"] = numberAt(publishedRow, modelSDField)
			row["published_live_sd"] = numberAt(publishedRow, liveSDField)
		}
		rows = append(rows, row)
	}
	if len(rows) != 5 {
		return nil, fmt.Errorf("%s: expected 5 published comparison rows, got %d", c.Name, len(rows))
	}
	return map[string]any{
		"source":      stringAt(reported, "source"),
		"metric":      stringAt(reported, "metric"),
		"model_field": c.PublishedModelField,
		"rows":        rows,
		"statistics":  comparisonStatistics(rows),
	}, nil
}

func markdownSummary(name string, analytics map[string]any) string {
	p, _ := analytics["procrustes"].(map[string]any)
	b, _ := analytics["bering"].(map[string]any)
	var out strings.Builder
	fmt.Fprintf(&out, "# %s reference summary\n\n", name)
	fmt.Fprintf(&out, "- Published comparison mode: `%s`.\n", stringAt(analytics, "mode"))
	fmt.Fprintf(&out, "- Procrustes goal status: `%s`; assessed services: %.0f; blockers: %.0f.\n", stringAt(p, "goal_status"), asFloat(p["services_assessed"]), asFloat(p["blockers"]))
	fmt.Fprintf(&out, "- Bering model: %.0f services, %.0f edges (%.0f async), %.0f endpoints.\n", asFloat(b["services"]), asFloat(b["edges"]), asFloat(b["async_edges"]), asFloat(b["endpoints"]))
	if publication, ok := analytics["publication"].(map[string]any); ok {
		fmt.Fprintf(&out, "- Published source: %s, DOI `%s`.\n", stringAt(publication, "title"), stringAt(publication, "doi"))
	}
	out.WriteString("\n| p_fail | k | Sheaft | Published model | Δ model | Published live | Δ live |\n")
	out.WriteString("|---:|---:|---:|---:|---:|---:|---:|\n")
	if comparison, ok := analytics["published_comparison"].(map[string]any); ok {
		for _, raw := range arrayAt(comparison, "rows") {
			row, ok := raw.(map[string]any)
			if ok {
				fmt.Fprintf(&out, "| %.1f | %.0f | %.6f | %.6f | %+.6f | %.6f | %+.6f |\n",
					numberAt(row, "failure_fraction"), numberAt(row, "fixed_k_failures"), numberAt(row, "sheaft"),
					numberAt(row, "published_model"), numberAt(row, "delta_from_published_model"),
					numberAt(row, "published_live"), numberAt(row, "delta_from_published_live"))
			}
		}
	}
	return out.String()
}

func profilesByName(report map[string]any) map[string]map[string]any {
	out := map[string]map[string]any{}
	for _, raw := range arrayAt(report, "profiles") {
		profile, ok := raw.(map[string]any)
		if ok {
			out[stringAt(profile, "name")] = profile
		}
	}
	return out
}

func comparisonStatistics(rows []any) map[string]any {
	return map[string]any{
		"against_published_model": seriesComparisonStatistics(rows, "published_model"),
		"against_published_live":  seriesComparisonStatistics(rows, "published_live"),
	}
}

func seriesComparisonStatistics(rows []any, referenceField string) map[string]any {
	actuals := make([]float64, 0, len(rows))
	references := make([]float64, 0, len(rows))
	meanSigned, meanAbsolute, meanSquare, maxAbsolute := 0.0, 0.0, 0.0, 0.0
	for _, raw := range rows {
		row, _ := raw.(map[string]any)
		actual := numberAt(row, "sheaft")
		reference := numberAt(row, referenceField)
		delta := actual - reference
		absolute := math.Abs(delta)
		actuals = append(actuals, actual)
		references = append(references, reference)
		meanSigned += delta
		meanAbsolute += absolute
		meanSquare += delta * delta
		if absolute > maxAbsolute {
			maxAbsolute = absolute
		}
	}
	n := float64(len(rows))
	return map[string]any{
		"mean_signed_error":      meanSigned / n,
		"mean_absolute_error":    meanAbsolute / n,
		"root_mean_square_error": math.Sqrt(meanSquare / n),
		"maximum_absolute_error": maxAbsolute,
		"pearson_correlation":    pearsonCorrelation(actuals, references),
	}
}

func pearsonCorrelation(left, right []float64) any {
	if len(left) == 0 || len(left) != len(right) {
		return nil
	}
	leftMean, rightMean := 0.0, 0.0
	for idx := range left {
		leftMean += left[idx]
		rightMean += right[idx]
	}
	leftMean /= float64(len(left))
	rightMean /= float64(len(right))
	numerator, leftSquare, rightSquare := 0.0, 0.0, 0.0
	for idx := range left {
		leftDelta := left[idx] - leftMean
		rightDelta := right[idx] - rightMean
		numerator += leftDelta * rightDelta
		leftSquare += leftDelta * leftDelta
		rightSquare += rightDelta * rightDelta
	}
	if leftSquare == 0 || rightSquare == 0 {
		return nil
	}
	return numerator / math.Sqrt(leftSquare*rightSquare)
}

func writeAggregateComparison(workRoot string, cases []caseSpec) error {
	families := map[string]any{}
	analyticsByCase := map[string]map[string]any{}
	for _, c := range cases {
		var analytics map[string]any
		if err := readJSON(filepath.Join(workRoot, c.Name, "analytics.json"), &analytics); err != nil {
			return fmt.Errorf("read %s analytics for aggregate comparison: %w", c.Name, err)
		}
		analyticsByCase[c.Name] = analytics
		family, _ := families[c.Family].(map[string]any)
		if family == nil {
			family = map[string]any{
				"publication": analytics["publication"],
				"variants":    map[string]any{},
			}
			families[c.Family] = family
		}
		variants, _ := family["variants"].(map[string]any)
		variants[c.Mode] = analytics["published_comparison"]
	}
	aggregate := map[string]any{
		"schema_version": "1.0.0",
		"description":    "Full Sheaft comparison against every aggregate result row reported by both published experiments.",
		"families":       families,
	}
	if err := writeJSON(filepath.Join(workRoot, "comparison.json"), aggregate); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(workRoot, "FULL_COMPARISON.md"), []byte(aggregateComparisonMarkdown(analyticsByCase)), 0o644)
}

func aggregateComparisonMarkdown(analyticsByCase map[string]map[string]any) string {
	var out strings.Builder
	out.WriteString("# Full comparison with published results\n\n")
	out.WriteString("These tables cover every aggregate row reported in Table 1 of the Social Network article and Table 2 of the OpenTelemetry Demo article. `Δ` is Sheaft minus the published value. Published live values are archival context; no live or chaos experiment is rerun here.\n\n")
	out.WriteString("## Social Network — ICSE-NIER 2026\n\n")
	out.WriteString("| Mode | p_fail | Sheaft | Published model | Δ model | Published live | Δ live |\n")
	out.WriteString("|---|---:|---:|---:|---:|---:|---:|\n")
	for _, mode := range []string{"norepl", "repl"} {
		analytics := analyticsByCase["social-network-"+mode]
		comparison, _ := analytics["published_comparison"].(map[string]any)
		for _, raw := range arrayAt(comparison, "rows") {
			row, _ := raw.(map[string]any)
			fmt.Fprintf(&out, "| %s | %.1f | %.6f | %.6f | %+.6f | %.6f | %+.6f |\n",
				mode, numberAt(row, "failure_fraction"), numberAt(row, "sheaft"), numberAt(row, "published_model"),
				numberAt(row, "delta_from_published_model"), numberAt(row, "published_live"), numberAt(row, "delta_from_published_live"))
		}
	}
	out.WriteString("\n## OpenTelemetry Demo — AINA 2026\n\n")
	out.WriteString("| p_fail | Sheaft all-block | Published all-block | Δ all-block | Sheaft async | Published async | Δ async | Published live |\n")
	out.WriteString("|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	allRows := comparisonRowsByFailureFraction(analyticsByCase["otel-demo-all-blocking"])
	asyncRows := comparisonRowsByFailureFraction(analyticsByCase["otel-demo-async"])
	for _, p := range []float64{0.1, 0.3, 0.5, 0.7, 0.9} {
		key := fmt.Sprintf("%.1f", p)
		all := allRows[key]
		async := asyncRows[key]
		fmt.Fprintf(&out, "| %.1f | %.6f | %.6f | %+.6f | %.6f | %.6f | %+.6f | %.6f |\n",
			p, numberAt(all, "sheaft"), numberAt(all, "published_model"), numberAt(all, "delta_from_published_model"),
			numberAt(async, "sheaft"), numberAt(async, "published_model"), numberAt(async, "delta_from_published_model"), numberAt(all, "published_live"))
	}
	out.WriteString("\nMachine-readable row deltas and aggregate error/correlation statistics are in `comparison.json`.\n")
	return out.String()
}

func comparisonRowsByFailureFraction(analytics map[string]any) map[string]map[string]any {
	out := map[string]map[string]any{}
	comparison, _ := analytics["published_comparison"].(map[string]any)
	for _, raw := range arrayAt(comparison, "rows") {
		row, _ := raw.(map[string]any)
		out[fmt.Sprintf("%.1f", numberAt(row, "failure_fraction"))] = row
	}
	return out
}

func semanticHashes(root string) (map[string]string, error) {
	return semanticHashesFor(root, referenceFiles)
}

func semanticHashesFor(root string, files []string) (map[string]string, error) {
	result := map[string]string{}
	for _, rel := range files {
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
		for key, raw := range typed {
			text, ok := raw.(string)
			if !ok {
				continue
			}
			switch {
			case (key == "source_ref" || key == "artifact_id") && strings.HasPrefix(text, "bering://discover?input="):
				typed[key] = normalizedDiscoveryRef(text)
			case key == "ref" && strings.HasSuffix(strings.ReplaceAll(text, "\\", "/"), "/procrustes/bering.overlay.yaml"):
				typed[key] = "paper://procrustes/bering.overlay.yaml"
			case key == "topology_version" && strings.HasPrefix(text, "sha256:"):
				typed[key] = zeroDigest()
			case key == "snapshot_id" && strings.HasPrefix(text, "snap-"):
				typed[key] = "snap-000000000000000000000000"
			case parent == "input_artifact" && key == "digest" && strings.HasPrefix(text, "sha256:"):
				typed[key] = zeroDigest()
			}
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

func normalizedDiscoveryRef(value string) string {
	lower := strings.ToLower(value)
	switch {
	case strings.Contains(lower, "topology-all-blocking.yaml"):
		return "bering://discover?input=case://otel-demo/topology-all-blocking.yaml"
	case strings.Contains(lower, "topology-norepl.yaml"):
		return "bering://discover?input=case://social-network/topology-norepl.yaml"
	case strings.Contains(lower, "/otel-demo/") || strings.Contains(lower, "%2fotel-demo%2f"):
		return "bering://discover?input=case://otel-demo/topology-api.yaml"
	case strings.Contains(lower, "/social-network/") || strings.Contains(lower, "%2fsocial-network%2f"):
		return "bering://discover?input=case://social-network/topology-api.yaml"
	default:
		return "bering://discover?input=case://unknown/topology-api.yaml"
	}
}

func zeroDigest() string {
	return "sha256:" + strings.Repeat("0", 64)
}

func writeReferenceCase(work, destination string) error {
	if err := os.RemoveAll(destination); err != nil {
		return err
	}
	return writeReferenceFiles(work, destination, referenceFiles)
}

func writeReferenceFiles(work, destination string, files []string) error {
	for _, rel := range files {
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
	if configured := strings.TrimSpace(os.Getenv("REPRODUCTION_ROOT")); configured != "" {
		root, err := filepath.Abs(configured)
		if err != nil {
			return "", fmt.Errorf("resolve REPRODUCTION_ROOT: %w", err)
		}
		if _, err := os.Stat(filepath.Join(root, "toolchain-lock.json")); err != nil {
			return "", fmt.Errorf("REPRODUCTION_ROOT %s does not contain toolchain-lock.json", root)
		}
		return root, nil
	}
	_, file, _, ok := runtime.Caller(0)
	if ok {
		root := filepath.Dir(file)
		if _, err := os.Stat(filepath.Join(root, "toolchain-lock.json")); err == nil {
			return root, nil
		}
	}
	workingDirectory, err := os.Getwd()
	if err == nil {
		if _, statErr := os.Stat(filepath.Join(workingDirectory, "toolchain-lock.json")); statErr == nil {
			return workingDirectory, nil
		}
	}
	return "", errors.New("cannot locate reproduction root; set REPRODUCTION_ROOT")
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
