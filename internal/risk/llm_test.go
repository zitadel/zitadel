package risk

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOllamaClientClassify(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"response":"{\"classification\":\"high\",\"confidence\":0.92,\"reason\":\"impossible context change\"}"}`)
	}))
	defer server.Close()

	client := NewOllamaClient(LLMConfig{
		Mode:               LLMModeObserve,
		Endpoint:           server.URL,
		Model:              "phi3:mini",
		Timeout:            time.Second,
		MaxEvents:          4,
		HighRiskConfidence: 0.85,
	}, server.Client())

	classification, err := client.Classify(context.Background(), Prompt{System: "test-system", User: "hello"})
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if classification.Normalized() != "high" {
		t.Fatalf("unexpected classification: %s", classification.Classification)
	}
	if classification.Confidence != 0.92 {
		t.Fatalf("unexpected confidence: %f", classification.Confidence)
	}
}

func TestOllamaClientClassify_TruncatedJSON(t *testing.T) {
	t.Parallel()

	// Simulate a model that truncates mid-string due to num_predict limit.
	truncated := `{"response":"{\"classification\":\"high\",\"confidence\":1.0,\"reason\":\"The anomaly is a change in the IP address from 172.18.0.1 to 172.18.0.2, indicating a potential session hijacking or session fixation attack. This is a high-risk anomaly as it could lead to unauthorized access to the system or data. The change in the fingerprint ID suggests a successful"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, truncated)
	}))
	defer server.Close()

	client := NewOllamaClient(LLMConfig{
		Mode:      LLMModeObserve,
		Endpoint:  server.URL,
		Model:     "test",
		Timeout:   time.Second,
		MaxEvents: 4,
	}, server.Client())

	classification, err := client.Classify(context.Background(), Prompt{System: "sys", User: "ctx"})
	if err != nil {
		t.Fatalf("expected repair to succeed, got: %v", err)
	}
	if classification.Normalized() != "high" {
		t.Errorf("classification = %q, want %q", classification.Normalized(), "high")
	}
}

func TestOllamaClientClassify_ConfidenceScale100(t *testing.T) {
	t.Parallel()

	// Some models return confidence on a 0-100 scale; should be normalized to 0-1.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"response":"{\"classification\":\"high\",\"confidence\":100,\"reason\":\"ip changed\"}"}`)
	}))
	defer server.Close()

	client := NewOllamaClient(LLMConfig{
		Mode:      LLMModeObserve,
		Endpoint:  server.URL,
		Model:     "test",
		Timeout:   time.Second,
		MaxEvents: 4,
	}, server.Client())

	c, err := client.Classify(context.Background(), Prompt{System: "sys", User: "ctx"})
	if err != nil {
		t.Fatalf("expected normalization to succeed, got: %v", err)
	}
	if c.Confidence != 1.0 {
		t.Errorf("confidence = %f, want 1.0 (normalized from 100)", c.Confidence)
	}
}

func TestRepairTruncatedJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantOut string
	}{
		{
			name:    "already complete",
			input:   `{"classification":"high","confidence":0.9,"reason":"ok"}`,
			wantOut: `{"classification":"high","confidence":0.9,"reason":"ok"}`,
		},
		{
			name:    "truncated inside string",
			input:   `{"classification":"high","confidence":1.0,"reason":"session hijack`,
			wantOut: `{"classification":"high","confidence":1.0,"reason":"session hijack"}`,
		},
		{
			name:    "truncated after comma (outside string)",
			input:   `{"classification":"high","confidence":1.0,`,
			wantOut: `{"classification":"high","confidence":1.0,}`,
		},
		{
			name:    "not a JSON object",
			input:   `"just a string"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repairTruncatedJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("repairTruncatedJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.wantOut {
				t.Errorf("repairTruncatedJSON() = %q, want %q", got, tt.wantOut)
			}
		})
	}
}
