package risk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LLMClient interface {
	Classify(ctx context.Context, prompt Prompt) (Classification, error)
}

type Classification struct {
	Classification string  `json:"classification"`
	Confidence     float64 `json:"confidence"`
	Reason         string  `json:"reason"`
}

func (c Classification) Normalized() string {
	return strings.ToLower(strings.TrimSpace(c.Classification))
}

func (c Classification) HighRisk() bool {
	switch c.Normalized() {
	case "high", "malicious", "block":
		return true
	default:
		return false
	}
}

// ollamaFormat is the legacy string format identifier for unconstrained JSON mode.
// This is significantly faster than JSON Schema constrained decoding on CPU because
// it avoids llama.cpp grammar-based token filtering.
var ollamaFormat = json.RawMessage(`"json"`)

type OllamaClient struct {
	endpoint    string
	model       string
	numPredict  int
	temperature *float64
	topK        int
	topP        float64
	keepAlive   string
	httpClient  *http.Client
}

func NewOllamaClient(cfg LLMConfig, httpClient *http.Client) *OllamaClient {
	client := httpClient
	if client == nil {
		client = &http.Client{}
	}
	cloned := *client
	cloned.Timeout = cfg.Timeout
	keepAlive := cfg.KeepAlive
	if keepAlive == "" {
		keepAlive = "10m"
	}
	return &OllamaClient{
		endpoint:    strings.TrimRight(cfg.Endpoint, "/"),
		model:       cfg.Model,
		numPredict:  cfg.NumPredict,
		temperature: cfg.Temperature,
		topK:        cfg.TopK,
		topP:        cfg.TopP,
		keepAlive:   keepAlive,
		httpClient:  &cloned,
	}
}

// ollamaOptions maps to Ollama's per-request model options.
type ollamaOptions struct {
	NumPredict  int      `json:"num_predict,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	TopK        int      `json:"top_k,omitempty"`
	TopP        float64  `json:"top_p,omitempty"`
}

type ollamaGenerateRequest struct {
	Model     string          `json:"model"`
	System    string          `json:"system,omitempty"`
	Prompt    string          `json:"prompt"`
	Format    json.RawMessage `json:"format"`
	Options   ollamaOptions   `json:"options,omitempty"`
	Stream    bool            `json:"stream"`
	KeepAlive string          `json:"keep_alive,omitempty"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}

func (c *OllamaClient) Classify(ctx context.Context, prompt Prompt) (Classification, error) {
	reqBody, err := json.Marshal(ollamaGenerateRequest{
		Model:  c.model,
		System: prompt.System,
		Prompt: prompt.User,
		Format: ollamaFormat,
		Options: ollamaOptions{
			NumPredict:  c.numPredict,
			Temperature: c.temperature,
			TopK:        c.topK,
			TopP:        c.topP,
		},
		Stream:    false,
		KeepAlive: c.keepAlive,
	})
	if err != nil {
		return Classification{}, fmt.Errorf("marshal ollama request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return Classification{}, fmt.Errorf("build ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Classification{}, fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Classification{}, fmt.Errorf("read ollama response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Classification{}, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var generateResp ollamaGenerateResponse
	if err := json.Unmarshal(body, &generateResp); err != nil {
		return Classification{}, fmt.Errorf("decode ollama response envelope: %w", err)
	}
	if strings.TrimSpace(generateResp.Error) != "" {
		return Classification{}, fmt.Errorf("ollama generate error: %s", strings.TrimSpace(generateResp.Error))
	}

	var classification Classification
	if err := json.Unmarshal([]byte(generateResp.Response), &classification); err != nil {
		return Classification{}, fmt.Errorf("decode ollama classification: %w (raw: %q)", err, generateResp.Response)
	}
	if classification.Normalized() == "" {
		return Classification{}, fmt.Errorf("ollama classification must not be empty")
	}
	if classification.Confidence < 0 || classification.Confidence > 1 {
		return Classification{}, fmt.Errorf("ollama confidence must be between 0 and 1, got %f", classification.Confidence)
	}
	return classification, nil
}
