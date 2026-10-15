package main

import (
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
