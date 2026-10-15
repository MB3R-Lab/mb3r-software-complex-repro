package main

import "testing"

func TestSuccessRequiresApplicationResult(t *testing.T) {
	for _, tc := range []struct {
		status int
		body   string
		want   bool
	}{
		{200, `[{"id":"p1"}]`, true}, {200, `[]`, false}, {200, `{}`, false},
		{200, `[{"name":"empty id"}]`, false}, {500, `[{"id":"p1"}]`, false},
	} {
		if got := success(tc.status, []byte(tc.body)); got != tc.want {
			t.Errorf("status=%d body=%s: %v", tc.status, tc.body, got)
		}
	}
}
