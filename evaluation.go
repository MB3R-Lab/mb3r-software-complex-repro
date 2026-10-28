package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const decisionEndpoint = "frontend:GET /api/recommendations"

type numericOracleCheck struct {
	ID            string  `json:"id"`
	Capability    string  `json:"capability"`
	Expected      float64 `json:"expected"`
	Observed      float64 `json:"observed"`
	AbsoluteError float64 `json:"absolute_error"`
	Tolerance     float64 `json:"tolerance"`
	Status        string  `json:"status"`
}

type processMeasurement struct {
	Label          string `json:"label"`
	DurationMillis int64  `json:"duration_ms"`
	PeakRSSKiB     *int64 `json:"peak_rss_kib,omitempty"`
}

func runDecisionWalkthrough(root, work string, bins map[string]string, toolchainVersion string) error {
	if err := os.RemoveAll(work); err != nil {
		return err
	}
	if err := os.MkdirAll(work, 0o755); err != nil {
		return err
	}
	caseDir := filepath.Join(root, "cases", "otel-demo")
	fixtureDir := filepath.Join(root, "cases", "decision-walkthrough")
	preflightDir := filepath.Join(work, "procrustes")
	if err := os.MkdirAll(preflightDir, 0o755); err != nil {
		return err
	}
	if err := runCommand("decision walkthrough Procrustes", caseDir, bins["procrustes"],
		"analyze", "--source", "paper="+filepath.Join(caseDir, "static"),
		"--goal", filepath.Join(caseDir, "goal.yaml"), "--out-dir", preflightDir, "--fail-on-blocked"); err != nil {
		return err
	}

	models := map[string]string{}
	for _, variant := range []string{"baseline", "regressed", "remediated"} {
		variantDir := filepath.Join(work, variant, "bering")
		if err := os.MkdirAll(variantDir, 0o755); err != nil {
			return err
		}
		modelPath := filepath.Join(variantDir, "model.json")
		if err := runCommand("decision walkthrough Bering "+variant, caseDir, bins["bering"],
			"discover", "--input", filepath.Join(caseDir, "topology-api.yaml"),
			"--overlay", filepath.Join(preflightDir, "bering.overlay.yaml"),
			"--replicas", filepath.Join(fixtureDir, "replicas-"+variant+".yaml"),
			"--discovered-at", discoveredAt,
			"--out", modelPath,
			"--snapshot-out", filepath.Join(variantDir, "snapshot.json"),
			"--quality-out", filepath.Join(variantDir, "model-quality.json"),
			"--snapshot-quality-out", filepath.Join(variantDir, "snapshot-quality.json")); err != nil {
			return err
		}
		models[variant] = modelPath
	}

	baselineAnalysis := filepath.Join(fixtureDir, "analysis-baseline.yaml")
	baselineOut := filepath.Join(work, "baseline", "sheaft")
	if err := runSheaftWithExpectedExit("decision baseline", bins["sheaft"], models["baseline"], baselineAnalysis, baselineOut, 0); err != nil {
		return err
	}

	candidateAnalysis := filepath.Join(work, "analysis-candidate.yaml")
	if err := renderTemplate(filepath.Join(fixtureDir, "analysis-candidate.template.yaml"), candidateAnalysis,
		map[string]string{"__BASELINE_MODEL__": yamlString(models["baseline"])}); err != nil {
		return err
	}
	candidateOut := filepath.Join(work, "regressed", "sheaft")
	if err := runSheaftWithExpectedExit("decision regressed candidate", bins["sheaft"], models["regressed"], candidateAnalysis, candidateOut, 2); err != nil {
		return err
	}

	importanceOut := filepath.Join(work, "regressed", "criticality")
	if err := runSheaftWithExpectedExit("decision candidate criticality", bins["sheaft"], models["regressed"], filepath.Join(fixtureDir, "importance.yaml"), importanceOut, 0); err != nil {
		return err
	}

	remediatedOut := filepath.Join(work, "remediated", "sheaft")
	if err := runSheaftWithExpectedExit("decision remediated candidate", bins["sheaft"], models["remediated"], candidateAnalysis, remediatedOut, 0); err != nil {
		return err
	}

	var baseline, candidate, remediated, importance map[string]any
	for path, destination := range map[string]*map[string]any{
		filepath.Join(baselineOut, "report.json"):   &baseline,
		filepath.Join(candidateOut, "report.json"):  &candidate,
		filepath.Join(remediatedOut, "report.json"): &remediated,
		filepath.Join(importanceOut, "report.json"): &importance,
	} {
		if err := readJSON(path, destination); err != nil {
			return err
		}
	}

	baselineTolerance, err := reportCertifiedTolerance(baseline, decisionEndpoint)
	if err != nil {
		return fmt.Errorf("decision baseline: %w", err)
	}
	candidateTolerance, err := reportCertifiedTolerance(candidate, decisionEndpoint)
	if err != nil {
		return fmt.Errorf("decision candidate: %w", err)
	}
	remediatedTolerance, err := reportCertifiedTolerance(remediated, decisionEndpoint)
	if err != nil {
		return fmt.Errorf("decision remediation: %w", err)
	}
	baselineAvailability, err := reportEndpointAvailability(baseline, "release-candidate", decisionEndpoint)
	if err != nil {
		return err
	}
	candidateAvailability, err := reportEndpointAvailability(candidate, "release-candidate", decisionEndpoint)
	if err != nil {
		return err
	}
	remediatedAvailability, err := reportEndpointAvailability(remediated, "release-candidate", decisionEndpoint)
	if err != nil {
		return err
	}
	recommendation, err := reportImportanceComponent(importance, decisionEndpoint, "service", "recommendation")
	if err != nil {
		return err
	}
	ad, err := reportImportanceComponent(importance, decisionEndpoint, "service", "ad")
	if err != nil {
		return err
	}

	if got := reportDecision(baseline); got != "pass" {
		return fmt.Errorf("decision baseline decision = %q, want pass", got)
	}
	if got := reportDecision(candidate); got != "fail" {
		return fmt.Errorf("decision candidate decision = %q, want fail", got)
	}
	if got := reportDecision(remediated); got != "pass" {
		return fmt.Errorf("decision remediation decision = %q, want pass", got)
	}
	if candidateTolerance >= baselineTolerance {
		return fmt.Errorf("decision candidate tolerance %.4f did not regress from baseline %.4f", candidateTolerance, baselineTolerance)
	}
	if remediatedTolerance < baselineTolerance {
		return fmt.Errorf("decision remediation tolerance %.4f did not recover baseline %.4f", remediatedTolerance, baselineTolerance)
	}
	if numberAt(recommendation, "birnbaum_importance") <= numberAt(ad, "birnbaum_importance") {
		return errors.New("decision criticality did not rank recommendation above ad")
	}
	reasonIDs := reportReasonIDs(candidate)
	if !containsString(reasonIDs, "boundary_regressed") {
		return fmt.Errorf("decision candidate reasons %v do not contain boundary_regressed", reasonIDs)
	}

	result := map[string]any{
		"schema_version":    "1.0.0",
		"toolchain_release": toolchainVersion,
		"scenario":          "slo-regression-and-importance-guided-remediation",
		"application":       "OpenTelemetry Demo",
		"endpoint":          decisionEndpoint,
		"slo":               0.95,
		"sampling": map[string]any{
			"distribution":       "independent replica fail-stop",
			"failure_population": "failure-eligible replicas in the Bering model",
			"seed":               2701,
			"trials":             20000,
			"importance_trials":  12000,
		},
		"changes": []any{
			map[string]any{"stage": "baseline", "recommendation_replicas": 2, "ad_replicas": 2},
			map[string]any{"stage": "regressed", "recommendation_replicas": 1, "ad_replicas": 2},
			map[string]any{"stage": "remediated", "recommendation_replicas": 2, "ad_replicas": 1},
		},
		"stages": []any{
			decisionStage("baseline", baseline, baselineAvailability, baselineTolerance),
			decisionStage("regressed", candidate, candidateAvailability, candidateTolerance),
			decisionStage("remediated", remediated, remediatedAvailability, remediatedTolerance),
		},
		"candidate_importance": map[string]any{
			"interpretation": "Birnbaum is used to nominate a component for inspection; the remediated model is evaluated independently.",
			"recommendation": recommendation,
			"ad":             ad,
		},
		"outcome": "REGRESSION_DETECTED_AND_REMEDIATION_VERIFIED",
	}
	if err := writeJSON(filepath.Join(work, "result.json"), result); err != nil {
		return err
	}
	if err := compareEvaluationReference(filepath.Join(root, "reference", "evaluation", "decision.json"), filepath.Join(work, "result.json"), toolchainVersion); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(work, "summary.md"), []byte(decisionMarkdown(result)), 0o644); err != nil {
		return err
	}
	fmt.Printf("decision-walkthrough-ok toolchain=%s endpoint=%q tolerance=%.2f->%.2f->%.2f decisions=pass,fail,pass\n",
		toolchainVersion, decisionEndpoint, baselineTolerance, candidateTolerance, remediatedTolerance)
	return nil
}

func runValidationMatrix(root, work string, bins map[string]string, toolchainVersion string) error {
	if err := os.RemoveAll(work); err != nil {
		return err
	}
	if err := os.MkdirAll(work, 0o755); err != nil {
		return err
	}
	checks := make([]any, 0, 16)
	procrustesChecks, err := runProcrustesValidation(root, filepath.Join(work, "procrustes"), bins["procrustes"])
	if err != nil {
		return err
	}
	checks = append(checks, procrustesChecks...)
	beringChecks, err := runBeringValidation(root, filepath.Join(work, "bering"), bins["bering"])
	if err != nil {
		return err
	}
	checks = append(checks, beringChecks...)
	oracleChecks, err := runSheaftOracleValidation(root, filepath.Join(work, "sheaft-oracles"), bins)
	if err != nil {
		return err
	}
	for _, check := range oracleChecks {
		checks = append(checks, check)
	}
	contractDir := filepath.Join(work, "contract-rejection")
	if err := runSheaftContractWalkthrough(root, contractDir, bins["sheaft"]); err != nil {
		return err
	}
	checks = append(checks, map[string]any{
		"id": "contract.unsupported-bering-schema", "component": "Procrustes-Bering-Sheaft boundary",
		"expected": "UNSUPPORTED_CONTRACT_REJECTED", "observed": "UNSUPPORTED_CONTRACT_REJECTED", "status": "PASS",
	})

	result := map[string]any{
		"schema_version":    "1.0.0",
		"toolchain_release": toolchainVersion,
		"scenario":          "controlled-validation-matrix",
		"checks":            checks,
		"summary": map[string]any{
			"total":  len(checks),
			"passed": len(checks),
			"failed": 0,
		},
		"stochastic_protocol": map[string]any{
			"distribution":       "independent replica fail-stop with perfect edges",
			"failure_population": "two failure-eligible oracle services; gateway excluded",
			"trials":             50000,
			"seed":               2711,
		},
	}
	if err := writeJSON(filepath.Join(work, "validation-matrix.json"), result); err != nil {
		return err
	}
	if err := compareEvaluationReference(filepath.Join(root, "reference", "evaluation", "validation-matrix.json"), filepath.Join(work, "validation-matrix.json"), toolchainVersion); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(work, "validation-matrix.md"), []byte(validationMarkdown(checks)), 0o644); err != nil {
		return err
	}
	fmt.Printf("validation-matrix-ok toolchain=%s checks=%d passed=%d\n", toolchainVersion, len(checks), len(checks))
	return nil
}

func runProcrustesValidation(root, work, binary string) ([]any, error) {
	fixture := filepath.Join(root, "cases", "validation", "procrustes")
	ready := filepath.Join(fixture, "ready")
	goal := filepath.Join(fixture, "goal.yaml")
	checks := make([]any, 0, 4)

	run := func(name, source, goalPath, wantStatus, wantRule string, wantExit int) error {
		out := filepath.Join(work, name)
		if err := os.MkdirAll(out, 0o755); err != nil {
			return err
		}
		exitCode, output, err := runCommandWithExit(fixture, binary, "analyze", "--source", "fixture="+source, "--goal", goalPath, "--out-dir", out, "--fail-on-blocked")
		if err != nil {
			return err
		}
		if exitCode != wantExit {
			return fmt.Errorf("Procrustes validation %s exit = %d, want %d\n%s", name, exitCode, wantExit, output)
		}
		observed := "INPUT_REJECTED"
		rules := []string{}
		if wantExit != 1 {
			var report map[string]any
			if err := readJSON(filepath.Join(out, "report.json"), &report); err != nil {
				return err
			}
			observed = stringAt(report, "goal_status")
			rules = append(rules, findingRules(arrayAt(report, "blockers"))...)
			rules = append(rules, findingRules(arrayAt(report, "remediation"))...)
			if observed != wantStatus {
				return fmt.Errorf("Procrustes validation %s status = %q, want %q", name, observed, wantStatus)
			}
			if wantRule != "" && !containsString(rules, wantRule) {
				return fmt.Errorf("Procrustes validation %s rules %v do not contain %s", name, rules, wantRule)
			}
		} else {
			if !strings.Contains(strings.ToLower(output), "unsupported") {
				return fmt.Errorf("Procrustes validation %s did not explain the unsupported input\n%s", name, output)
			}
			if _, err := os.Stat(filepath.Join(out, "report.json")); !os.IsNotExist(err) {
				return fmt.Errorf("Procrustes validation %s unexpectedly produced report.json", name)
			}
		}
		checks = append(checks, map[string]any{
			"id": "procrustes." + name, "component": "Procrustes", "expected": wantStatus,
			"observed": observed, "expected_exit_code": wantExit, "observed_exit_code": exitCode,
			"required_rule": wantRule, "observed_rules": uniqueSorted(rules), "status": "PASS",
		})
		return nil
	}

	if err := run("ready", ready, goal, "EXPECTED", "", 0); err != nil {
		return nil, err
	}
	conflictSource := filepath.Join(work, "inputs", "identity-conflict")
	if err := copyFixtureWithIdentityConflict(ready, conflictSource); err != nil {
		return nil, err
	}
	if err := run("identity-conflict", conflictSource, goal, "BLOCKED", "identity.unique", 2); err != nil {
		return nil, err
	}
	missingSource := filepath.Join(work, "inputs", "missing-scenario")
	if err := os.MkdirAll(missingSource, 0o755); err != nil {
		return nil, err
	}
	manifest, err := os.ReadFile(filepath.Join(ready, "manifests.yaml"))
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(missingSource, "manifests.yaml"), manifest, 0o644); err != nil {
		return nil, err
	}
	if err := run("missing-scenario", missingSource, goal, "BLOCKED", "operation.scenario", 2); err != nil {
		return nil, err
	}
	if err := run("unsupported-goal-contract", ready, filepath.Join(fixture, "unsupported-goal.yaml"), "INPUT_REJECTED", "", 1); err != nil {
		return nil, err
	}
	return checks, nil
}

func runBeringValidation(root, work, binary string) ([]any, error) {
	fixture := filepath.Join(root, "cases", "validation", "bering")
	if err := os.MkdirAll(work, 0o755); err != nil {
		return nil, err
	}
	modelPath := filepath.Join(work, "model.json")
	if err := runCommand("Bering controlled trace validation", fixture, binary,
		"discover", "--input", filepath.Join(fixture, "operations.json"), "--discovered-at", discoveredAt,
		"--out", modelPath, "--snapshot-out", filepath.Join(work, "snapshot.json"),
		"--quality-out", filepath.Join(work, "model-quality.json"),
		"--snapshot-quality-out", filepath.Join(work, "snapshot-quality.json")); err != nil {
		return nil, err
	}
	var model, expected map[string]any
	if err := readJSON(modelPath, &model); err != nil {
		return nil, err
	}
	if err := readJSON(filepath.Join(fixture, "expected.json"), &expected); err != nil {
		return nil, err
	}
	checks := make([]any, 0, 3)
	for _, spec := range []struct {
		name string
		rows []any
	}{
		{"services", arrayAt(model, "services")},
		{"edges", arrayAt(model, "edges")},
		{"endpoints", arrayAt(model, "endpoints")},
	} {
		actual := idsFromRows(spec.rows)
		want := stringsFromAny(expected[spec.name])
		missing, extra := setDifference(want, actual), setDifference(actual, want)
		if len(missing) > 0 || len(extra) > 0 {
			return nil, fmt.Errorf("Bering %s mismatch: missing=%v extra=%v", spec.name, missing, extra)
		}
		checks = append(checks, map[string]any{
			"id": "bering.trace-" + spec.name, "component": "Bering", "expected": want, "observed": actual,
			"true_positives": len(want), "missing": missing, "extra": extra, "status": "PASS",
		})
	}
	return checks, nil
}

func runSheaftOracleValidation(root, work string, bins map[string]string) ([]numericOracleCheck, error) {
	fixture := filepath.Join(root, "cases", "validation", "oracle")
	beringOut := filepath.Join(work, "bering")
	if err := os.MkdirAll(beringOut, 0o755); err != nil {
		return nil, err
	}
	modelPath := filepath.Join(beringOut, "model.json")
	if err := runCommand("oracle Bering model", fixture, bins["bering"],
		"discover", "--input", filepath.Join(fixture, "topology.yaml"), "--discovered-at", discoveredAt,
		"--out", modelPath, "--snapshot-out", filepath.Join(beringOut, "snapshot.json"),
		"--quality-out", filepath.Join(beringOut, "model-quality.json"),
		"--snapshot-quality-out", filepath.Join(beringOut, "snapshot-quality.json")); err != nil {
		return nil, err
	}

	outputs := map[string]map[string]any{}
	exitCodes := map[string]int{}
	for _, gateCase := range []struct {
		name      string
		minimum   string
		wantExit  int
		wantCause string
	}{
		{"pass", "0.05", 0, "gate_pass"},
		{"fail", "0.20", 2, "boundary_below_minimum"},
		{"indeterminate", "0.10", 2, "boundary_indeterminate"},
	} {
		analysis := filepath.Join(work, "analysis-"+gateCase.name+".yaml")
		if err := renderTemplate(filepath.Join(fixture, "analysis.template.yaml"), analysis,
			map[string]string{"__MINIMUM_TOLERANCE__": gateCase.minimum}); err != nil {
			return nil, err
		}
		out := filepath.Join(work, gateCase.name)
		if err := runSheaftWithExpectedExit("oracle "+gateCase.name, bins["sheaft"], modelPath, analysis, out, gateCase.wantExit); err != nil {
			return nil, err
		}
		var report map[string]any
		if err := readJSON(filepath.Join(out, "report.json"), &report); err != nil {
			return nil, err
		}
		if !containsString(reportReasonIDs(report), gateCase.wantCause) {
			return nil, fmt.Errorf("oracle %s reasons %v do not contain %s", gateCase.name, reportReasonIDs(report), gateCase.wantCause)
		}
		outputs[gateCase.name] = report
		exitCodes[gateCase.name] = gateCase.wantExit
	}

	pass := outputs["pass"]
	immediate, err := reportEndpointAvailability(pass, "oracle-state", "gateway:GET /immediate")
	if err != nil {
		return nil, err
	}
	eventual, err := reportEndpointAvailability(pass, "oracle-state", "gateway:POST /eventual")
	if err != nil {
		return nil, err
	}
	component, err := reportImportanceComponent(pass, "gateway:GET /immediate", "service", "sync-worker")
	if err != nil {
		return nil, err
	}
	tolerance, err := reportCertifiedTolerance(pass, "gateway:GET /immediate")
	if err != nil {
		return nil, err
	}
	checks := []numericOracleCheck{
		newNumericCheck("sheaft.immediate-availability", "immediate response excludes async work", 0.90, immediate, 0.012),
		newNumericCheck("sheaft.eventual-availability", "eventual completion requires sync and async work", 0.81, eventual, 0.012),
		newNumericCheck("sheaft.birnbaum-series-service", "paired Birnbaum importance", 1.0, numberAt(component, "birnbaum_importance"), 0.001),
		newNumericCheck("sheaft.certified-slo-boundary", "confidence-certified evaluated-grid tolerance", 0.05, tolerance, 0.0),
		newNumericCheck("sheaft.gate-pass", "SLO gate exit code", 0.0, float64(exitCodes["pass"]), 0.0),
		newNumericCheck("sheaft.gate-fail", "SLO gate exit code", 2.0, float64(exitCodes["fail"]), 0.0),
		newNumericCheck("sheaft.gate-indeterminate", "SLO gate exit code", 2.0, float64(exitCodes["indeterminate"]), 0.0),
	}
	for _, check := range checks {
		if check.Status != "PASS" {
			return nil, fmt.Errorf("oracle %s: observed %.6f, expected %.6f ± %.6f", check.ID, check.Observed, check.Expected, check.Tolerance)
		}
	}
	return checks, nil
}

func runCostProfile(root, work string, bins map[string]string, toolchainVersion string, repeats int) error {
	if err := os.RemoveAll(work); err != nil {
		return err
	}
	if err := os.MkdirAll(work, 0o755); err != nil {
		return err
	}
	caseDir := filepath.Join(root, "cases", "otel-demo")
	fixtureDir := filepath.Join(root, "cases", "decision-walkthrough")
	full := make([]processMeasurement, 0, repeats)
	reanalysis := make([]processMeasurement, 0, repeats)
	for iteration := 1; iteration <= repeats; iteration++ {
		iterationDir := filepath.Join(work, fmt.Sprintf("run-%02d", iteration))
		preflight := filepath.Join(iterationDir, "procrustes")
		bering := filepath.Join(iterationDir, "bering")
		analysis := filepath.Join(iterationDir, "analysis")
		for _, dir := range []string{preflight, bering, analysis} {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
		}
		started := time.Now()
		measurements := make([]processMeasurement, 0, 3)
		measurement, err := measureCommand("procrustes", caseDir, bins["procrustes"],
			"analyze", "--source", "paper="+filepath.Join(caseDir, "static"),
			"--goal", filepath.Join(caseDir, "goal.yaml"), "--out-dir", preflight, "--fail-on-blocked")
		if err != nil {
			return err
		}
		measurements = append(measurements, measurement)
		modelPath := filepath.Join(bering, "model.json")
		measurement, err = measureCommand("bering", caseDir, bins["bering"],
			"discover", "--input", filepath.Join(caseDir, "topology-api.yaml"),
			"--overlay", filepath.Join(preflight, "bering.overlay.yaml"),
			"--replicas", filepath.Join(fixtureDir, "replicas-baseline.yaml"), "--discovered-at", discoveredAt,
			"--out", modelPath, "--snapshot-out", filepath.Join(bering, "snapshot.json"),
			"--quality-out", filepath.Join(bering, "model-quality.json"),
			"--snapshot-quality-out", filepath.Join(bering, "snapshot-quality.json"))
		if err != nil {
			return err
		}
		measurements = append(measurements, measurement)
		measurement, err = measureCommand("sheaft", caseDir, bins["sheaft"],
			"run", "--model", modelPath, "--analysis", filepath.Join(fixtureDir, "analysis-baseline.yaml"), "--out-dir", analysis)
		if err != nil {
			return err
		}
		measurements = append(measurements, measurement)
		peak := maximumRSS(measurements)
		full = append(full, processMeasurement{Label: "full-pipeline", DurationMillis: time.Since(started).Milliseconds(), PeakRSSKiB: peak})

		reanalysisDir := filepath.Join(iterationDir, "reanalysis")
		measurement, err = measureCommand("saved-model-analysis", caseDir, bins["sheaft"],
			"run", "--model", modelPath, "--analysis", filepath.Join(fixtureDir, "analysis-baseline.yaml"), "--out-dir", reanalysisDir)
		if err != nil {
			return err
		}
		reanalysis = append(reanalysis, measurement)
	}
	result := map[string]any{
		"schema_version":    "1.0.0",
		"toolchain_release": toolchainVersion,
		"scenario":          "descriptive-application-cost",
		"repeats":           repeats,
		"environment": map[string]any{
			"os": runtime.GOOS, "architecture": runtime.GOARCH, "logical_cpus": runtime.NumCPU(),
			"memory_metric": "maximum resident set size of any child process; Linux reports KiB",
		},
		"full_pipeline":        map[string]any{"runs": full, "summary": measurementSummary(full)},
		"saved_model_analysis": map[string]any{"runs": reanalysis, "summary": measurementSummary(reanalysis)},
		"interpretation":       "Descriptive cost on the stated machine; not a cross-system performance benchmark.",
	}
	if err := writeJSON(filepath.Join(work, "cost-profile.json"), result); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(work, "cost-profile.md"), []byte(costMarkdown(result)), 0o644); err != nil {
		return err
	}
	fullSummary, _ := result["full_pipeline"].(map[string]any)
	fullStats, _ := fullSummary["summary"].(map[string]any)
	reSummary, _ := result["saved_model_analysis"].(map[string]any)
	reStats, _ := reSummary["summary"].(map[string]any)
	fmt.Printf("cost-profile-ok toolchain=%s repeats=%d full_median_ms=%.0f reanalysis_median_ms=%.0f\n",
		toolchainVersion, repeats, asFloat(fullStats["duration_median_ms"]), asFloat(reStats["duration_median_ms"]))
	return nil
}

func runSheaftWithExpectedExit(label, binary, model, analysis, out string, wantExit int) error {
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	exitCode, output, err := runCommandWithExit(filepath.Dir(analysis), binary,
		"run", "--model", model, "--analysis", analysis, "--out-dir", out, "--why")
	if err != nil {
		return fmt.Errorf("%s: %w", label, err)
	}
	if exitCode != wantExit {
		return fmt.Errorf("%s exit code %d, want %d\n%s", label, exitCode, wantExit, output)
	}
	return nil
}

func renderTemplate(source, destination string, replacements map[string]string) error {
	raw, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	text := string(raw)
	for placeholder, replacement := range replacements {
		if !strings.Contains(text, placeholder) {
			return fmt.Errorf("template %s does not contain %s", source, placeholder)
		}
		text = strings.ReplaceAll(text, placeholder, replacement)
	}
	if strings.Contains(text, "__") {
		return fmt.Errorf("template %s still contains an unresolved placeholder", source)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destination, []byte(text), 0o644)
}

func compareEvaluationReference(reference, generated, currentRelease string) error {
	want, err := normalizedFile(reference)
	if err != nil {
		return fmt.Errorf("read evaluation reference %s: %w", reference, err)
	}
	got, err := normalizedFile(generated)
	if err != nil {
		return fmt.Errorf("read generated evaluation result %s: %w", generated, err)
	}
	var historical, current map[string]any
	historicalDecoder := json.NewDecoder(bytes.NewReader(want))
	historicalDecoder.UseNumber()
	if err := historicalDecoder.Decode(&historical); err != nil {
		return fmt.Errorf("decode historical evaluation reference: %w", err)
	}
	currentDecoder := json.NewDecoder(bytes.NewReader(got))
	currentDecoder.UseNumber()
	if err := currentDecoder.Decode(&current); err != nil {
		return fmt.Errorf("decode generated evaluation result: %w", err)
	}
	historicalRelease, _ := historical["toolchain_release"].(string)
	generatedRelease, _ := current["toolchain_release"].(string)
	if historicalRelease == "" || generatedRelease != currentRelease {
		return fmt.Errorf("evaluation release provenance mismatch: historical=%q generated=%q current=%q", historicalRelease, generatedRelease, currentRelease)
	}
	// Only release provenance differs between a new toolchain run and an
	// immutable historical reference. Compare every other normalized field.
	current["toolchain_release"] = historicalRelease
	got, err = json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}
	want, err = json.MarshalIndent(historical, "", "  ")
	if err != nil {
		return err
	}
	if !bytes.Equal(want, got) {
		return fmt.Errorf("evaluation reference mismatch: generated %s differs from %s", generated, reference)
	}
	return nil
}

func yamlString(value string) string {
	return strconv.Quote(filepath.ToSlash(value))
}

func reportDecision(report map[string]any) string {
	evaluation, _ := report["policy_evaluation"].(map[string]any)
	return stringAt(evaluation, "decision")
}

func reportReasonIDs(report map[string]any) []string {
	evaluation, _ := report["policy_evaluation"].(map[string]any)
	reasons := arrayAt(evaluation, "reasons")
	ids := make([]string, 0, len(reasons))
	for _, raw := range reasons {
		if reason, ok := raw.(map[string]any); ok {
			ids = append(ids, stringAt(reason, "id"))
		}
	}
	return uniqueSorted(ids)
}

func reportEndpointAvailability(report map[string]any, profileName, endpointID string) (float64, error) {
	for _, raw := range arrayAt(report, "profiles") {
		profile, _ := raw.(map[string]any)
		if stringAt(profile, "name") != profileName {
			continue
		}
		simulation, _ := profile["simulation"].(map[string]any)
		endpoints, _ := simulation["endpoint_availability"].(map[string]any)
		value, ok := endpoints[endpointID]
		if !ok {
			return 0, fmt.Errorf("profile %s has no endpoint %s", profileName, endpointID)
		}
		return asFloat(value), nil
	}
	return 0, fmt.Errorf("report has no profile %s", profileName)
}

func reportCertifiedTolerance(report map[string]any, endpointID string) (float64, error) {
	for _, sweepRaw := range arrayAt(report, "sweeps") {
		sweep, _ := sweepRaw.(map[string]any)
		for _, boundaryRaw := range arrayAt(sweep, "boundaries") {
			boundary, _ := boundaryRaw.(map[string]any)
			if stringAt(boundary, "endpoint_id") != endpointID {
				continue
			}
			certified, _ := boundary["certified_tolerance"].(map[string]any)
			if certified == nil {
				return 0, fmt.Errorf("endpoint %s has no certified tolerance", endpointID)
			}
			return numberAt(certified, "axis_value"), nil
		}
	}
	return 0, fmt.Errorf("report has no sweep boundary for %s", endpointID)
}

func reportImportanceComponent(report map[string]any, endpointID, kind, id string) (map[string]any, error) {
	for _, importanceRaw := range arrayAt(report, "importance") {
		importance, _ := importanceRaw.(map[string]any)
		for _, endpointRaw := range arrayAt(importance, "endpoints") {
			endpoint, _ := endpointRaw.(map[string]any)
			if stringAt(endpoint, "endpoint_id") != endpointID {
				continue
			}
			for _, componentRaw := range arrayAt(endpoint, "components") {
				component, _ := componentRaw.(map[string]any)
				if stringAt(component, "kind") == kind && stringAt(component, "id") == id {
					return component, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("importance report has no %s/%s for %s", kind, id, endpointID)
}

func decisionStage(name string, report map[string]any, availability, tolerance float64) map[string]any {
	return map[string]any{
		"name": name, "decision": reportDecision(report), "availability": availability,
		"certified_tolerance": tolerance, "reason_ids": reportReasonIDs(report),
	}
}

func newNumericCheck(id, capability string, expected, observed, tolerance float64) numericOracleCheck {
	errorValue := math.Abs(observed - expected)
	status := "PASS"
	if errorValue > tolerance+1e-12 {
		status = "FAIL"
	}
	return numericOracleCheck{ID: id, Capability: capability, Expected: expected, Observed: observed, AbsoluteError: errorValue, Tolerance: tolerance, Status: status}
}

func findingRules(findings []any) []string {
	rules := make([]string, 0, len(findings))
	for _, raw := range findings {
		if finding, ok := raw.(map[string]any); ok {
			if rule := stringAt(finding, "rule_id"); rule != "" {
				rules = append(rules, rule)
			}
		}
	}
	return rules
}

func copyFixtureWithIdentityConflict(source, destination string) error {
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}
	for _, name := range []string{"checkout.yaml", "manifests.yaml"} {
		raw, err := os.ReadFile(filepath.Join(source, name))
		if err != nil {
			return err
		}
		if name == "manifests.yaml" {
			before := "service.name=$(OTEL_SERVICE_NAME),service.version=1.2.3"
			after := "service.name=wrong,service.version=1.2.3"
			if !bytes.Contains(raw, []byte(before)) {
				return errors.New("identity-conflict fixture source no longer contains the expected field")
			}
			raw = bytes.Replace(raw, []byte(before), []byte(after), 1)
		}
		if err := os.WriteFile(filepath.Join(destination, name), raw, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func idsFromRows(rows []any) []string {
	ids := make([]string, 0, len(rows))
	for _, raw := range rows {
		if row, ok := raw.(map[string]any); ok {
			ids = append(ids, stringAt(row, "id"))
		}
	}
	return uniqueSorted(ids)
}

func stringsFromAny(value any) []string {
	values, _ := value.([]any)
	out := make([]string, 0, len(values))
	for _, raw := range values {
		if text, ok := raw.(string); ok {
			out = append(out, text)
		}
	}
	return uniqueSorted(out)
}

func setDifference(left, right []string) []string {
	rightSet := make(map[string]struct{}, len(right))
	for _, value := range right {
		rightSet[value] = struct{}{}
	}
	result := make([]string, 0)
	for _, value := range left {
		if _, ok := rightSet[value]; !ok {
			result = append(result, value)
		}
	}
	return result
}

func uniqueSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func measureCommand(label, dir, name string, args ...string) (processMeasurement, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	started := time.Now()
	err := cmd.Run()
	measurement := processMeasurement{Label: label, DurationMillis: time.Since(started).Milliseconds()}
	if cmd.ProcessState != nil {
		measurement.PeakRSSKiB = peakRSSKiB(cmd.ProcessState.SysUsage())
	}
	if err != nil {
		return measurement, fmt.Errorf("measure %s: %w\n%s", label, err, strings.TrimSpace(output.String()))
	}
	return measurement, nil
}

func peakRSSKiB(usage any) *int64 {
	if usage == nil {
		return nil
	}
	value := reflect.ValueOf(usage)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}
	field := value.FieldByName("Maxrss")
	if !field.IsValid() {
		return nil
	}
	var rss int64
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rss = field.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rss = int64(field.Uint())
	default:
		return nil
	}
	if runtime.GOOS == "darwin" {
		rss /= 1024
	}
	return &rss
}

func maximumRSS(measurements []processMeasurement) *int64 {
	var maximum *int64
	for _, measurement := range measurements {
		if measurement.PeakRSSKiB == nil {
			continue
		}
		if maximum == nil || *measurement.PeakRSSKiB > *maximum {
			value := *measurement.PeakRSSKiB
			maximum = &value
		}
	}
	return maximum
}

func measurementSummary(measurements []processMeasurement) map[string]any {
	durations := make([]int64, 0, len(measurements))
	var peak *int64
	for _, measurement := range measurements {
		durations = append(durations, measurement.DurationMillis)
		if measurement.PeakRSSKiB != nil && (peak == nil || *measurement.PeakRSSKiB > *peak) {
			value := *measurement.PeakRSSKiB
			peak = &value
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	summary := map[string]any{
		"duration_min_ms": durations[0], "duration_median_ms": medianInt64(durations), "duration_max_ms": durations[len(durations)-1],
	}
	if peak != nil {
		summary["peak_rss_kib"] = *peak
	}
	return summary
}

func medianInt64(sorted []int64) float64 {
	middle := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return float64(sorted[middle])
	}
	return float64(sorted[middle-1]+sorted[middle]) / 2
}

func decisionMarkdown(result map[string]any) string {
	var out strings.Builder
	out.WriteString("# SLO regression and remediation walkthrough\n\n")
	fmt.Fprintf(&out, "Application: OpenTelemetry Demo. Endpoint: `%s`. SLO: `%.2f`.\n\n", stringAt(result, "endpoint"), numberAt(result, "slo"))
	out.WriteString("| Stage | Decision | Availability | Certified failure tolerance |\n|---|---|---:|---:|\n")
	for _, raw := range arrayAt(result, "stages") {
		stage, _ := raw.(map[string]any)
		fmt.Fprintf(&out, "| %s | %s | %.4f | %.2f |\n", stringAt(stage, "name"), stringAt(stage, "decision"), numberAt(stage, "availability"), numberAt(stage, "certified_tolerance"))
	}
	out.WriteString("\nBirnbaum importance nominates a component for inspection; only the independent re-analysis of the changed model establishes the reported effect.\n")
	return out.String()
}

func validationMarkdown(checks []any) string {
	var out strings.Builder
	out.WriteString("# Controlled validation matrix\n\n")
	out.WriteString("| Check | Component | Expected | Observed | Error / mismatch | Result |\n|---|---|---|---|---|---|\n")
	for _, raw := range checks {
		encoded, _ := json.Marshal(raw)
		var row map[string]any
		_ = json.Unmarshal(encoded, &row)
		expected := fmt.Sprint(row["expected"])
		observed := fmt.Sprint(row["observed"])
		errorText := "—"
		if value, ok := row["absolute_error"]; ok {
			errorText = fmt.Sprintf("%v (tol %v)", value, row["tolerance"])
		} else if missing, ok := row["missing"]; ok {
			errorText = fmt.Sprintf("missing=%v; extra=%v", missing, row["extra"])
		}
		component := fmt.Sprint(row["component"])
		if component == "<nil>" {
			component = fmt.Sprint(row["capability"])
		}
		fmt.Fprintf(&out, "| `%s` | %s | %s | %s | %s | %s |\n",
			fmt.Sprint(row["id"]), component, expected, observed, errorText, fmt.Sprint(row["status"]))
	}
	return out.String()
}

func costMarkdown(result map[string]any) string {
	full, _ := result["full_pipeline"].(map[string]any)
	fullSummary, _ := full["summary"].(map[string]any)
	reanalysis, _ := result["saved_model_analysis"].(map[string]any)
	reSummary, _ := reanalysis["summary"].(map[string]any)
	var out strings.Builder
	out.WriteString("# Descriptive application cost\n\n")
	out.WriteString("These measurements describe this machine and are not a cross-system performance benchmark.\n\n")
	out.WriteString("| Path | Repeats | Median time (ms) | Peak child RSS (KiB) |\n|---|---:|---:|---:|\n")
	fmt.Fprintf(&out, "| Full Procrustes → Bering → Sheaft pass | %.0f | %.0f | %s |\n",
		numberAt(result, "repeats"), asFloat(fullSummary["duration_median_ms"]), optionalNumber(fullSummary, "peak_rss_kib"))
	fmt.Fprintf(&out, "| Sheaft analysis of saved Bering model | %.0f | %.0f | %s |\n",
		numberAt(result, "repeats"), asFloat(reSummary["duration_median_ms"]), optionalNumber(reSummary, "peak_rss_kib"))
	return out.String()
}

func optionalNumber(value map[string]any, key string) string {
	if _, ok := value[key]; !ok {
		return "unavailable"
	}
	return fmt.Sprintf("%.0f", asFloat(value[key]))
}
