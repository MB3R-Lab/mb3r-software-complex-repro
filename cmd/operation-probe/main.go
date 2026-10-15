// The probe controls a disposable CI deployment and preserves observed evidence.
// It does not implement a resilience simulator or replace application handlers.
package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type channel struct {
	mu          sync.Mutex
	disabled    bool
	upstream    string
	connections map[net.Conn]bool
}

func (c *channel) forward(client net.Conn) {
	c.mu.Lock()
	if c.disabled {
		c.mu.Unlock()
		client.Close()
		return
	}
	c.connections[client] = true
	c.mu.Unlock()
	defer func() { client.Close(); c.mu.Lock(); delete(c.connections, client); c.mu.Unlock() }()
	server, err := net.DialTimeout("tcp", c.upstream, 2*time.Second)
	if err != nil {
		return
	}
	defer server.Close()
	c.mu.Lock()
	if c.disabled {
		c.mu.Unlock()
		return
	}
	c.connections[server] = true
	c.mu.Unlock()
	defer func() { c.mu.Lock(); delete(c.connections, server); c.mu.Unlock() }()
	done := make(chan struct{}, 1)
	go func() { io.Copy(server, client); done <- struct{}{} }()
	io.Copy(client, server)
	client.Close()
	server.Close()
	<-done
}

func proxy() error {
	channels := map[string]*channel{}
	for _, spec := range []struct{ name, address, upstream string }{
		{"frontend-recommendation", ":19001", "recommendation:9001"},
		{"recommendation-catalog", ":13550", "product-catalog:3550"},
		{"frontend-catalog", ":23550", "product-catalog:3550"},
	} {
		c := &channel{upstream: spec.upstream, connections: map[net.Conn]bool{}}
		channels[spec.name] = c
		listener, err := net.Listen("tcp", spec.address)
		if err != nil {
			return err
		}
		go func() {
			for {
				client, err := listener.Accept()
				if err != nil {
					return
				}
				go c.forward(client)
			}
		}()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", 405)
			return
		}
		var input struct {
			Disabled []string `json:"disabled"`
		}
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			http.Error(w, "invalid JSON", 400)
			return
		}
		disabled := map[string]bool{}
		for _, name := range input.Disabled {
			if channels[name] == nil {
				http.Error(w, "unknown channel", 400)
				return
			}
			disabled[name] = true
		}
		for name, c := range channels {
			c.mu.Lock()
			c.disabled = disabled[name]
			if c.disabled {
				for conn := range c.connections {
					conn.Close()
				}
			}
			c.mu.Unlock()
		}
		json.NewEncoder(w).Encode(input)
	})
	return http.ListenAndServe(":8474", mux)
}

type observation struct {
	Condition      string `json:"condition"`
	Attempt        int    `json:"attempt"`
	TraceID        string `json:"trace_id"`
	Status         int    `json:"http_status"`
	Body           string `json:"body"`
	TransportError string `json:"transport_error,omitempty"`
	Success        bool   `json:"success"`
	Expected       bool   `json:"expected_success"`
	DurationMS     int64  `json:"duration_ms"`
}

func success(status int, body []byte) bool {
	if status != 200 {
		return false
	}
	var products []struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(body, &products) != nil || len(products) == 0 {
		return false
	}
	for _, p := range products {
		if p.ID == "" {
			return false
		}
	}
	return true
}

func request(condition string, attempt int) observation {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		panic(err)
	}
	traceID := hex.EncodeToString(raw[:16])
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/api/recommendations?currencyCode=USD", nil)
	req.Header.Set("traceparent", "00-"+traceID+"-"+hex.EncodeToString(raw[16:])+"-01")
	o := observation{Condition: condition, Attempt: attempt, TraceID: traceID, Expected: condition == "healthy"}
	start := time.Now()
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	o.DurationMS = time.Since(start).Milliseconds()
	if err != nil {
		o.TransportError = err.Error()
		return o
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		o.TransportError = err.Error()
		return o
	}
	o.Status = resp.StatusCode
	o.Body = string(body)
	o.Success = success(o.Status, body)
	return o
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func setChannels(disabled []string) error {
	body, _ := json.Marshal(map[string]any{"disabled": disabled})
	resp, err := http.Post("http://127.0.0.1:8474/state", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("proxy status %d", resp.StatusCode)
	}
	return nil
}

func compose(args ...string) error {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func waitHealthy() error {
	for i := 0; i < 60; i++ {
		if request("recovery", i).Success {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("application did not recover to a successful nonempty response")
}

type jaegerTag struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
type jaegerTrace struct {
	Spans []struct {
		TraceID    string      `json:"traceID"`
		SpanID     string      `json:"spanID"`
		Operation  string      `json:"operationName"`
		ProcessID  string      `json:"processID"`
		Start      int64       `json:"startTime"`
		Duration   int64       `json:"duration"`
		Tags       []jaegerTag `json:"tags"`
		References []struct {
			Type   string `json:"refType"`
			SpanID string `json:"spanID"`
		} `json:"references"`
	} `json:"spans"`
	Processes map[string]struct {
		Service string `json:"serviceName"`
	} `json:"processes"`
}

func getTrace(id string) (json.RawMessage, []map[string]any, error) {
	var last []byte
	// Wait for exporters and Jaeger storage. Keep the actual unmodified response.
	time.Sleep(3 * time.Second)
	for attempt := 0; attempt < 15; attempt++ {
		resp, err := http.Get("http://127.0.0.1:16686/api/traces/" + id)
		if err == nil {
			last, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			var data struct {
				Data []jaegerTrace `json:"data"`
			}
			if json.Unmarshal(last, &data) == nil && len(data.Data) > 0 && len(data.Data[0].Spans) > 0 {
				normalized := []map[string]any{}
				for _, trace := range data.Data {
					for _, span := range trace.Spans {
						attrs := map[string]any{}
						kind := "internal"
						for _, tag := range span.Tags {
							attrs[tag.Key] = tag.Value
							if tag.Key == "span.kind" {
								kind = fmt.Sprint(tag.Value)
							}
						}
						parent := ""
						for _, ref := range span.References {
							if ref.Type == "CHILD_OF" {
								parent = ref.SpanID
								break
							}
						}
						normalized = append(normalized, map[string]any{"trace_id": span.TraceID, "span_id": span.SpanID, "parent_span_id": parent,
							"service": trace.Processes[span.ProcessID].Service, "name": span.Operation, "kind": kind, "attributes": attrs,
							"start_time": time.UnixMicro(span.Start).UTC().Format(time.RFC3339Nano), "end_time": time.UnixMicro(span.Start + span.Duration).UTC().Format(time.RFC3339Nano)})
					}
				}
				return json.RawMessage(last), normalized, nil
			}
		}
		time.Sleep(2 * time.Second)
	}
	return nil, nil, fmt.Errorf("no recorded trace %s: %s", id, last)
}

func collect() error {
	out := ".tmp/operation-evidence"
	if err := os.MkdirAll(out, 0755); err != nil {
		return err
	}
	if err := waitHealthy(); err != nil {
		return err
	}
	cases := []struct{ name, service, channel string }{
		{"healthy", "", ""}, {"recommendation-down", "recommendation", ""},
		{"catalog-down", "product-catalog", ""}, {"frontend-recommendation-down", "", "frontend-recommendation"},
		{"recommendation-catalog-down", "", "recommendation-catalog"}, {"frontend-catalog-down", "", "frontend-catalog"},
	}
	all := []observation{}
	healthySpans := []map[string]any{}
	for _, tc := range cases {
		if err := setChannels(nil); err != nil {
			return err
		}
		if err := waitHealthy(); err != nil {
			return err
		}
		if tc.service != "" {
			if err := compose("stop", "-t", "1", tc.service); err != nil {
				return err
			}
		}
		if tc.channel != "" {
			if err := setChannels([]string{tc.channel}); err != nil {
				return err
			}
		}
		state, err := exec.Command("docker", "compose", "ps", "--all", "--format", "json").Output()
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(out, tc.name+"-containers.json"), state, 0644); err != nil {
			return err
		}
		for n := 1; n <= 2; n++ {
			o := request(tc.name, n)
			all = append(all, o)
			if err := writeJSON(filepath.Join(out, "observations.json"), all); err != nil {
				return err
			}
			fmt.Printf("observe condition=%s attempt=%d status=%d success=%t transport_error=%q\n", tc.name, n, o.Status, o.Success, o.TransportError)
			raw, spans, err := getTrace(o.TraceID)
			if err != nil {
				return err
			}
			if err := writeJSON(filepath.Join(out, "traces", fmt.Sprintf("%s-%d.json", tc.name, n)), raw); err != nil {
				return err
			}
			if tc.name == "healthy" {
				healthySpans = append(healthySpans, spans...)
			}
			if o.TransportError != "" {
				return fmt.Errorf("inconclusive HTTP transport failure in %s", tc.name)
			}
			if o.Success != o.Expected {
				return fmt.Errorf("unexpected observed semantics in %s (response retained)", tc.name)
			}
		}
		if tc.service != "" {
			if err := compose("start", tc.service); err != nil {
				return err
			}
		}
	}
	if err := setChannels(nil); err != nil {
		return err
	}
	if err := waitHealthy(); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(out, "bering-traces.json"), map[string]any{"spans": healthySpans}); err != nil {
		return err
	}
	return writeJSON(filepath.Join(out, "collection.json"), map[string]any{"application_commit": os.Getenv("APP_REVISION"), "workflow_run": os.Getenv("GITHUB_RUN_ID"), "conditions": len(cases), "requests": len(all), "success_criterion": "HTTP 200 and nonempty product array with nonempty IDs", "status": "observed; model agreement not yet evaluated"})
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: operation-probe proxy|collect")
	}
	var err error
	switch os.Args[1] {
	case "proxy":
		err = proxy()
	case "collect":
		err = collect()
	default:
		err = fmt.Errorf("unknown mode")
	}
	if err != nil {
		log.Fatal(err)
	}
}
