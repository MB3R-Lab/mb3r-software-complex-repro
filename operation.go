package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func runOperationWalkthrough(root, work string, bins map[string]string, version string) error {
	fixture := filepath.Join(root, "cases", "operation-semantics")
	evidence := filepath.Join(fixture, "evidence")
	if err := os.MkdirAll(work, 0755); err != nil {
		return err
	}
	var plan, collection map[string]any
	if err := readJSON(filepath.Join(fixture, "plan.json"), &plan); err != nil {
		return err
	}
	if err := readJSON(filepath.Join(evidence, "collection.json"), &collection); err != nil {
		return err
	}
	if stringAt(plan, "application_commit") != stringAt(collection, "application_commit") {
		return fmt.Errorf("application revision mismatch")
	}
	observedAt, err := validateOperationEvidence(evidence)
	if err != nil {
		return err
	}
	preflight := filepath.Join(work, "procrustes")
	if err := os.MkdirAll(preflight, 0755); err != nil {
		return err
	}
	if err := runCommand("operation preflight", fixture, bins["procrustes"], "analyze", "--source", "execution="+filepath.Join(fixture, "static"), "--goal", filepath.Join(fixture, "goal.yaml"), "--out-dir", preflight, "--fail-on-blocked"); err != nil {
		return err
	}
	baseOverlay := filepath.Join(preflight, "bering.overlay.yaml")
	rawModel := filepath.Join(work, "observed-model.json")
	tracePath := filepath.Join(evidence, "bering-traces.json")
	if err := runCommand("operation discovery from observed traces", fixture, bins["bering"], "discover", "--input", tracePath, "--overlay", baseOverlay, "--discovered-at", observedAt, "--out", rawModel, "--snapshot-out", filepath.Join(work, "observed-snapshot.json")); err != nil {
		return err
	}
	var observed map[string]any
	if err := readJSON(rawModel, &observed); err != nil {
		return err
	}
	// Bind exact discovered edge IDs to the predeclared calls. No edge is removed.
	required := []string{}
	byPair := map[string]string{}
	for _, raw := range arrayAt(plan, "required_call_pairs") {
		items, ok := raw.([]any)
		if !ok || len(items) != 2 {
			return fmt.Errorf("invalid directed call pair")
		}
		pair := []string{fmt.Sprint(items[0]), fmt.Sprint(items[1])}
		matches := []string{}
		for _, edgeRaw := range arrayAt(observed, "edges") {
			edge := edgeRaw.(map[string]any)
			if stringAt(edge, "from") == pair[0] && stringAt(edge, "to") == pair[1] && stringAt(edge, "kind") == "sync" && edge["blocking"] == true {
				matches = append(matches, stringAt(edge, "id"))
			}
		}
		if len(matches) != 1 {
			return fmt.Errorf("expected one observed RPC channel %s -> %s, found %v; inspect traces before changing the declaration", pair[0], pair[1], matches)
		}
		required = append(required, matches[0])
		byPair[pair[0]+"->"+pair[1]] = matches[0]
	}
	sort.Strings(required)
	endpoint := stringAt(plan, "endpoint")
	predicateRef := stringAt(plan, "predicate_ref")
	predicate := map[string]any{"type": "edge_aware", "entry_service": "frontend", "mandatory_targets": []string{"recommendation", "product-catalog"}, "edge_modes": []string{"sync"}, "required_edges": required}
	predicatePath := filepath.Join(work, "predicate.json")
	if err := writeJSON(predicatePath, map[string]any{"schema_version": "1.1", "predicates": map[string]any{predicateRef: predicate}}); err != nil {
		return err
	}
	services := []any{}
	for _, raw := range arrayAt(observed, "services") {
		svc := raw.(map[string]any)
		id := stringAt(svc, "id")
		services = append(services, map[string]any{"id": id, "failure_eligible": id == "recommendation" || id == "product-catalog"})
	}
	overlayPath := filepath.Join(work, "operation-overlay.json")
	if err := writeJSON(overlayPath, map[string]any{"name": "observed-operation-required-calls", "services": services, "endpoints": []any{map[string]any{
		"id": endpoint, "entry_service": "frontend", "method": "GET", "path": "/api/recommendations", "predicate_ref": predicateRef,
		"semantics": map[string]any{"predicate_mode": "external_predicate", "mandatory_targets": []string{"recommendation", "product-catalog"}, "dependency_modes": []string{"sync"}, "source": "pinned-handler-and-recorded-fault-states"}}}}); err != nil {
		return err
	}
	models := map[string]string{}
	stages := []any{}
	var currentReport map[string]any
	for _, raw := range arrayAt(plan, "variants") {
		variant := raw.(map[string]any)
		name := stringAt(variant, "name")
		dir := filepath.Join(work, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		replicas := filepath.Join(dir, "replicas.json")
		if err := writeJSON(replicas, map[string]any{"recommendation": numberAt(variant, "recommendation"), "product-catalog": numberAt(variant, "product-catalog")}); err != nil {
			return err
		}
		modelPath := filepath.Join(dir, "model.json")
		models[name] = modelPath
		whatif := filepath.Join(dir, "replica-overlay.json")
		if err := writeJSON(whatif, map[string]any{"name": "declared-replica-whatif", "services": []any{
			map[string]any{"id": "recommendation", "replicas": variant["recommendation"]},
			map[string]any{"id": "product-catalog", "replicas": variant["product-catalog"]},
		}}); err != nil {
			return err
		}
		if err := runCommand("operation Bering "+name, fixture, bins["bering"], "discover", "--input", tracePath, "--overlay", baseOverlay, "--overlay", overlayPath, "--overlay", whatif, "--discovered-at", observedAt, "--out", modelPath, "--snapshot-out", filepath.Join(dir, "snapshot.json")); err != nil {
			return err
		}
		// Verify discovery actually applied the replica what-if, after all overlays.
		var actual map[string]any
		if err := readJSON(modelPath, &actual); err != nil {
			return err
		}
		for _, svcRaw := range arrayAt(actual, "services") {
			svc := svcRaw.(map[string]any)
			id := stringAt(svc, "id")
			if id == "recommendation" || id == "product-catalog" {
				if numberAt(svc, "replicas") != numberAt(variant, id) {
					return fmt.Errorf("replica override lost for %s/%s", name, id)
				}
			}
		}
		cfg := operationAnalysis(plan, predicatePath)
		if name != "baseline" {
			cfg["baselines"] = []any{map[string]any{"name": "accepted", "path": models["baseline"]}}
			gate := cfg["gate"].(map[string]any)
			rule := gate["boundary_rules"].([]any)[0].(map[string]any)
			rule["baseline"] = "accepted"
			rule["max_regression"] = plan["max_regression"]
		}
		if name == "current" {
			cfg["importance"] = []any{map[string]any{"name": "candidate-choice", "profile": "scenario", "method": "birnbaum", "confidence_level": 0.95, "targets": []any{map[string]any{"endpoint_id": endpoint}}, "components": map[string]any{"kinds": []string{"service", "edge"}}}}
		}
		configPath := filepath.Join(dir, "analysis.json")
		if err := writeJSON(configPath, cfg); err != nil {
			return err
		}
		report, err := runOperationAnalysis(bins["sheaft"], modelPath, configPath, filepath.Join(dir, "sheaft"))
		if err != nil {
			return err
		}
		if name == "current" {
			currentReport = report
		}
		availability, err := reportEndpointAvailability(report, "scenario", endpoint)
		if err != nil {
			return err
		}
		tolerance, err := reportCertifiedTolerance(report, endpoint)
		if err != nil {
			return err
		}
		theoretical := math.Pow(numberAt(plan, "edge_live_probability"), float64(len(required))) * (1 - math.Pow(numberAt(plan, "failure_probability"), numberAt(variant, "recommendation"))) * (1 - math.Pow(numberAt(plan, "failure_probability"), numberAt(variant, "product-catalog")))
		if math.Abs(availability-theoretical) > 0.006 {
			return fmt.Errorf("%s availability %v disagrees with independent series oracle %v", name, availability, theoretical)
		}
		row := decisionStage(name, report, availability, tolerance)
		row["replicas"] = variant
		row["analytical_availability"] = theoretical
		row["absolute_error"] = math.Abs(availability - theoretical)
		row["availability_wilson_95"] = operationWilson(availability, numberAt(plan, "trials"))
		stages = append(stages, row)
	}
	checks, err := checkOperationStates(plan, evidence, work, models["current"], predicatePath, observed, byPair, bins["sheaft"])
	if err != nil {
		return err
	}
	rec, err := reportImportanceComponent(currentReport, endpoint, "service", "recommendation")
	if err != nil {
		return err
	}
	cat, err := reportImportanceComponent(currentReport, endpoint, "service", "product-catalog")
	if err != nil {
		return err
	}
	result := map[string]any{"schema_version": "1.0.0", "toolchain_release": version, "endpoint": endpoint, "application_commit": plan["application_commit"], "evidence_workflow_run": collection["workflow_run"], "required_edges": required, "semantic_checks": checks, "stages": stages, "current_importance": map[string]any{"recommendation": rec, "product-catalog": cat}, "scope": plan["scope"], "plan": plan}
	if err := writeJSON(filepath.Join(work, "result.json"), result); err != nil {
		return err
	}
	var summary strings.Builder
	summary.WriteString("# Observed operation and candidate comparison\n\nEvery predeclared observed state agrees with the required-call predicate. Both catalog calls remain in the discovered graph.\n\n| Configuration | Availability | Analytical oracle | Grid tolerance | Decision | Reasons |\n|---|---:|---:|---:|---|---|\n")
	for _, raw := range stages {
		r := raw.(map[string]any)
		fmt.Fprintf(&summary, "| %s | %.6f | %.6f | %.2f | %s | %v |\n", stringAt(r, "name"), numberAt(r, "availability"), numberAt(r, "analytical_availability"), numberAt(r, "certified_tolerance"), stringAt(r, "decision"), r["reason_ids"])
	}
	summary.WriteString("\nEndpoint gates use point estimates. Grid certification uses pointwise Wilson intervals. Regression compares discrete certified grid metrics, not a confidence bound on the change of a continuous threshold.\n\nReplica changes are conditional model analyses; the collected requests validate the selected execution predicate, not replication failover or operational failure probabilities.\n")
	if err := os.WriteFile(filepath.Join(work, "summary.md"), []byte(summary.String()), 0644); err != nil {
		return err
	}
	fmt.Printf("operation-walkthrough-ok states=%d variants=%d endpoint=%q\n", len(checks), len(stages), endpoint)
	return nil
}

func operationAnalysis(plan map[string]any, predicate string) map[string]any {
	endpoint := stringAt(plan, "endpoint")
	return map[string]any{"schema_version": "1.3", "seed": plan["seed"], "trials": plan["trials"], "predicate_contract": predicate,
		"profiles": []any{map[string]any{"name": "scenario", "sampling_mode": "independent_replica", "failure_probability": plan["failure_probability"], "reliability": map[string]any{"edge_live_probability": plan["edge_live_probability"]}}},
		"sweeps":   []any{map[string]any{"name": "operation-boundary", "profile": "scenario", "confidence_level": 0.95, "axis": map[string]any{"type": "independent_replica_failure_probability", "values": plan["grid"]}, "targets": []any{map[string]any{"endpoint_id": endpoint, "slo": plan["slo"]}}}},
		"gate":     map[string]any{"mode": "fail", "default_action": "fail", "global_threshold": 0.0, "endpoint_thresholds": map[string]any{endpoint: plan["slo"]}, "boundary_rules": []any{map[string]any{"sweep": "operation-boundary", "endpoint_id": endpoint, "minimum_certified_tolerance": plan["minimum_certified_tolerance"]}}}}
}

func runOperationAnalysis(binary, model, config, out string) (map[string]any, error) {
	if err := os.MkdirAll(out, 0755); err != nil {
		return nil, err
	}
	code, output, err := runCommandWithExit(filepath.Dir(model), binary, "run", "--model", model, "--analysis", config, "--out-dir", out, "--why")
	if err != nil {
		return nil, err
	}
	if code != 0 && code != 2 {
		return nil, fmt.Errorf("analysis exit %d: %s", code, output)
	}
	var report map[string]any
	if err := readJSON(filepath.Join(out, "report.json"), &report); err != nil {
		return nil, err
	}
	return report, nil
}

func checkOperationStates(plan map[string]any, evidence, work, model, predicate string, observed map[string]any, pairs map[string]string, binary string) ([]any, error) {
	var observations []map[string]any
	if err := readJSON(filepath.Join(evidence, "observations.json"), &observations); err != nil {
		return nil, err
	}
	profiles := []any{}
	conditions := []string{"healthy", "recommendation-down", "catalog-down", "frontend-recommendation-down", "recommendation-catalog-down", "frontend-catalog-down"}
	for _, name := range conditions {
		svc := map[string]any{}
		for _, raw := range arrayAt(observed, "services") {
			svc[stringAt(raw.(map[string]any), "id")] = 1.0
		}
		edges := map[string]any{}
		switch name {
		case "recommendation-down":
			svc["recommendation"] = 0.0
		case "catalog-down":
			svc["product-catalog"] = 0.0
		case "frontend-recommendation-down":
			edges[pairs["frontend->recommendation"]] = 0.0
		case "recommendation-catalog-down":
			edges[pairs["recommendation->product-catalog"]] = 0.0
		case "frontend-catalog-down":
			edges[pairs["frontend->product-catalog"]] = 0.0
		}
		profiles = append(profiles, map[string]any{"name": name, "sampling_mode": "independent_replica", "reliability": map[string]any{"services": svc, "edges": edges, "edge_live_probability": 1.0}})
	}
	config := filepath.Join(work, "semantic-analysis.json")
	if err := writeJSON(config, map[string]any{"schema_version": "1.3", "seed": 2719, "trials": 1, "predicate_contract": predicate, "profiles": profiles, "gate": map[string]any{"mode": "report", "global_threshold": 0.0}}); err != nil {
		return nil, err
	}
	report, err := runOperationAnalysis(binary, model, config, filepath.Join(work, "semantic-check"))
	if err != nil {
		return nil, err
	}
	// Keep the same discovered graph and mandatory targets; change only the
	// predicate contract to reachability without required calls.
	legacyPredicate := filepath.Join(work, "reachability-only-predicate.json")
	if err := writeJSON(legacyPredicate, map[string]any{"schema_version": "1.0", "predicates": map[string]any{stringAt(plan, "predicate_ref"): map[string]any{
		"type": "edge_aware", "entry_service": "frontend", "mandatory_targets": []string{"recommendation", "product-catalog"}, "edge_modes": []string{"sync"},
	}}}); err != nil {
		return nil, err
	}
	legacyConfig := filepath.Join(work, "reachability-only-analysis.json")
	if err := writeJSON(legacyConfig, map[string]any{"schema_version": "1.3", "seed": 2719, "trials": 1, "predicate_contract": legacyPredicate, "profiles": profiles, "gate": map[string]any{"mode": "report", "global_threshold": 0.0}}); err != nil {
		return nil, err
	}
	legacy, err := runOperationAnalysis(binary, model, legacyConfig, filepath.Join(work, "reachability-only-check"))
	if err != nil {
		return nil, err
	}
	checks := []any{}
	for _, name := range conditions {
		value, err := reportEndpointAvailability(report, name, stringAt(plan, "endpoint"))
		if err != nil {
			return nil, err
		}
		count := 0
		for _, o := range observations {
			if stringAt(o, "condition") != name {
				continue
			}
			count++
			if stringAt(o, "transport_error") != "" {
				return nil, fmt.Errorf("inconclusive transport result for %s", name)
			}
			if o["success"] != (value == 1) {
				return nil, fmt.Errorf("predicate disagrees with observation for %s: model=%v response=%v", name, value, o)
			}
		}
		if count != 2 {
			return nil, fmt.Errorf("expected two recorded requests for %s, got %d", name, count)
		}
		oldValue, err := reportEndpointAvailability(legacy, name, stringAt(plan, "endpoint"))
		if err != nil {
			return nil, err
		}
		checks = append(checks, map[string]any{"condition": name, "observed_requests": count, "predicate": value == 1, "agreement": true, "reachability_only_predicate": oldValue == 1, "reachability_only_agreement": oldValue == value})
	}
	return checks, nil
}

// These intervals describe the ordinary profile estimate, not the sweep gate.
func operationWilson(p, n float64) map[string]float64 {
	z := 1.959963984540054
	d := 1 + z*z/n
	c := (p + z*z/(2*n)) / d
	h := z * math.Sqrt(p*(1-p)/n+z*z/(4*n*n)) / d
	return map[string]float64{"lower_bound": math.Max(0, c-h), "upper_bound": math.Min(1, c+h)}
}

func operationResponseSuccess(status float64, body string) bool {
	var products []struct {
		ID string `json:"id"`
	}
	if status != 200 || json.Unmarshal([]byte(body), &products) != nil || len(products) == 0 {
		return false
	}
	for _, product := range products {
		if strings.TrimSpace(product.ID) == "" {
			return false
		}
	}
	return true
}

// Validate response/trace linkage before invoking the tools. Use the last span
// end as a deterministic discovery timestamp, never an archival case timestamp.
func validateOperationEvidence(evidence string) (string, error) {
	var observations []map[string]any
	if err := readJSON(filepath.Join(evidence, "observations.json"), &observations); err != nil {
		return "", err
	}
	healthySpans := map[string]string{}
	seen := map[string]bool{}
	for _, o := range observations {
		id := stringAt(o, "trace_id")
		if id == "" || seen[id] || stringAt(o, "transport_error") != "" {
			return "", fmt.Errorf("invalid, duplicate, or inconclusive observation %q", id)
		}
		seen[id] = true
		if operationResponseSuccess(numberAt(o, "http_status"), stringAt(o, "body")) != o["success"] {
			return "", fmt.Errorf("stored success disagrees with actual response %s", id)
		}
		var raw map[string]any
		path := filepath.Join(evidence, "traces", fmt.Sprintf("%s-%.0f.json", stringAt(o, "condition"), numberAt(o, "attempt")))
		if err := readJSON(path, &raw); err != nil {
			return "", err
		}
		data := arrayAt(raw, "data")
		if len(data) != 1 || stringAt(data[0].(map[string]any), "traceID") != id {
			return "", fmt.Errorf("raw trace does not match response %s", id)
		}
		spans := arrayAt(data[0].(map[string]any), "spans")
		if len(spans) == 0 {
			return "", fmt.Errorf("empty trace %s", id)
		}
		if stringAt(o, "condition") == "healthy" {
			for _, s := range spans {
				span := s.(map[string]any)
				healthySpans[id+"/"+stringAt(span, "spanID")] = stringAt(span, "operationName")
			}
		}
	}
	var normalized map[string]any
	if err := readJSON(filepath.Join(evidence, "bering-traces.json"), &normalized); err != nil {
		return "", err
	}
	last := time.Time{}
	for _, s := range arrayAt(normalized, "spans") {
		span := s.(map[string]any)
		key := stringAt(span, "trace_id") + "/" + stringAt(span, "span_id")
		name, ok := healthySpans[key]
		if !ok || name != stringAt(span, "name") {
			return "", fmt.Errorf("normalized span has no matching recorded healthy span: %s", key)
		}
		delete(healthySpans, key)
		end, err := time.Parse(time.RFC3339Nano, stringAt(span, "end_time"))
		if err != nil {
			return "", err
		}
		if end.After(last) {
			last = end
		}
	}
	if len(healthySpans) != 0 || last.IsZero() {
		return "", fmt.Errorf("missing normalized healthy spans")
	}
	return last.UTC().Format(time.RFC3339Nano), nil
}
