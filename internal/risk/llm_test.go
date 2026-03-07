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
