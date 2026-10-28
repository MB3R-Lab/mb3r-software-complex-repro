#!/usr/bin/env python3
"""Actions-only counterfactual for an added database observation in one replay.

The original Bering output stays intact. This script writes separate diagnostic
models, removes only the single newly observed database node and its edge, and
reruns the same Sheaft configurations with file paths relocated to the runner.
"""

import argparse
import copy
import hashlib
import json
from pathlib import Path
import subprocess
import time


STAGES = ("baseline", "current", "add-recommendation", "add-catalog")


def read(path):
    return json.loads(path.read_text(encoding="utf-8"))


def write(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def digest(path):
    return "sha256:" + hashlib.sha256(path.read_bytes()).hexdigest()


def added_database(old, new):
    old_services = {item["id"] for item in old["services"]}
    new_services = {item["id"] for item in new["services"]}
    added = new_services - old_services
    if old_services - new_services or len(added) != 1:
        raise ValueError(f"expected exactly one added service: {sorted(added)}")
    database_id = next(iter(added))
    if not database_id.startswith("database:"):
        raise ValueError(f"added service is not a database: {database_id}")
    old_edges = {item["id"] for item in old["edges"]}
    new_edges = {item["id"] for item in new["edges"]}
    added_edges = [item for item in new["edges"] if item["id"] in new_edges - old_edges]
    if old_edges - new_edges or len(added_edges) != 1:
        raise ValueError(f"expected exactly one added edge: {added_edges}")
    edge = added_edges[0]
    if edge["from"] != "product-catalog" or edge["to"] != database_id:
        raise ValueError(f"new edge does not connect catalog to database: {edge}")
    return database_id, edge["id"]


def filter_model(model, database_id, edge_id):
    result = copy.deepcopy(model)
    services = result["services"]
    edges = result["edges"]
    if sum(item["id"] == database_id for item in services) != 1:
        raise ValueError("diagnostic model has no unique database node")
    if sum(item["id"] == edge_id for item in edges) != 1:
        raise ValueError("diagnostic model has no unique database edge")
    incident = [item for item in edges if item["from"] == database_id or item["to"] == database_id]
    if [item["id"] for item in incident] != [edge_id]:
        raise ValueError("database has additional incident edges")
    result["services"] = [item for item in services if item["id"] != database_id]
    result["edges"] = [item for item in edges if item["id"] != edge_id]
    return result


def endpoint_availability(report, endpoint):
    rows = [row for profile in report["profiles"] if profile["name"] == "scenario"
            for row in profile["endpoint_results"] if row["endpoint_id"] == endpoint]
    if len(rows) != 1:
        raise ValueError("expected one scenario endpoint result")
    return rows[0]["availability"]


def importance(report, endpoint, service):
    rows = [component for block in report["importance"]
            for row in block["endpoints"] if row["endpoint_id"] == endpoint
            for component in row["components"]
            if component["kind"] == "service" and component["id"] == service]
    if len(rows) != 1:
        raise ValueError(f"expected one importance component for {service}")
    return rows[0]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--work", required=True, type=Path)
    parser.add_argument("--reference", required=True, type=Path)
    parser.add_argument("--sheaft", required=True, type=Path)
    args = parser.parse_args()
    work = args.work.resolve()
    reference = args.reference.resolve()
    sheaft = args.sheaft.resolve()
    old_model = read(reference / "observed-model.json")
    new_model_path = work / "observed-model.json"
    database_id, edge_id = added_database(old_model, read(new_model_path))
    old_result = read(reference / "result.json")
    new_result_path = work / "result.json"
    new_result = read(new_result_path)
    endpoint = old_result["endpoint"]
    diag = work / "draw-order-diagnostic"
    diag.mkdir(parents=True, exist_ok=True)
    summary = {
        "criterion": "single newly observed database node and its sole edge removed only in diagnostic copies",
        "database_id": database_id,
        "edge_id": edge_id,
        "full_observed_model_sha256": digest(new_model_path),
        "full_result_sha256": digest(new_result_path),
        "historical_result_sha256": digest(reference / "result.json"),
        "stages": {},
    }
    started = time.monotonic()
    for stage in STAGES:
        original_dir = work / stage
        diagnostic_dir = diag / stage
        model_path = original_dir / "model.json"
        filtered_path = diagnostic_dir / "model.json"
        write(filtered_path, filter_model(read(model_path), database_id, edge_id))
        analysis = read(original_dir / "analysis.json")
        analysis["predicate_contract"] = str(work / "predicate.json")
        for baseline in analysis.get("baselines", []):
            if baseline.get("name") != "accepted":
                raise ValueError("unexpected baseline declaration")
            baseline["path"] = str(diag / "baseline" / "model.json")
        config_path = diagnostic_dir / "analysis.json"
        write(config_path, analysis)
        output = diagnostic_dir / "sheaft"
        start_stage = time.monotonic()
        run = subprocess.run([str(sheaft), "run", "--model", str(filtered_path),
                              "--analysis", str(config_path), "--out-dir", str(output), "--why"],
                             capture_output=True, text=True, check=False)
        if run.returncode not in (0, 2):
            raise RuntimeError(f"diagnostic Sheaft {stage} exit {run.returncode}: {run.stderr}\n{run.stdout}")
        report = read(output / "report.json")
        historical = next(row for row in old_result["stages"] if row["name"] == stage)
        full = next(row for row in new_result["stages"] if row["name"] == stage)
        value = endpoint_availability(report, endpoint)
        summary["stages"][stage] = {
            "full_model_sha256": digest(model_path),
            "diagnostic_model_sha256": digest(filtered_path),
            "historical_availability": historical["availability"],
            "full_availability": full["availability"],
            "diagnostic_availability": value,
            "exact_historical_availability": value == historical["availability"],
            "diagnostic_exit": run.returncode,
            "elapsed_seconds": time.monotonic() - start_stage,
        }
    report = read(diag / "current" / "sheaft" / "report.json")
    summary["importance"] = {}
    for service in ("recommendation", "product-catalog"):
        item = importance(report, endpoint, service)
        old = old_result["current_importance"][service]
        full = new_result["current_importance"][service]
        summary["importance"][service] = {
            "historical": old, "full": full, "diagnostic": item,
            "exact_historical": all(item.get(key) == value for key, value in old.items()),
        }
    summary["elapsed_seconds"] = time.monotonic() - started
    summary["draw_order_cause_supported"] = (
        all(row["exact_historical_availability"] for row in summary["stages"].values())
        and all(row["exact_historical"] for row in summary["importance"].values()))
    write(diag / "summary.json", summary)
    print(json.dumps({"draw_order_cause_supported": summary["draw_order_cause_supported"],
                      "stages": {name: row["exact_historical_availability"]
                                 for name, row in summary["stages"].items()},
                      "importance": {name: row["exact_historical"]
                                     for name, row in summary["importance"].items()},
                      "elapsed_seconds": summary["elapsed_seconds"]}, sort_keys=True))


if __name__ == "__main__":
    main()
