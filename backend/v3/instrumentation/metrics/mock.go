package metrics

import (
	"context"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// MockMetrics implements the metrics.Metrics interface for testing
type MockMetrics struct {
	mu              sync.RWMutex
	histogramValues map[string][]float64
	counterValues   map[string]int64
	histogramLabels map[string][]map[string]attribute.Value
	counterLabels   map[string][]map[string]attribute.Value
}

var _ Metrics = new(MockMetrics)

// NewMockMetrics creates a new Metrics instance for testing
func NewMockMetrics() *MockMetrics {
	return &MockMetrics{
		histogramValues: make(map[string][]float64),
		counterValues:   make(map[string]int64),
		histogramLabels: make(map[string][]map[string]attribute.Value),
		counterLabels:   make(map[string][]map[string]attribute.Value),
	}
}

func (m *MockMetrics) GetExporter() http.Handler {
	return nil
}

func (m *MockMetrics) GetMetricsProvider() metric.MeterProvider {
	return nil
}

func (m *MockMetrics) RegisterCounter(name, description string) error {
	return nil
}

func (m *MockMetrics) AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counterValues[name] += value
	m.counterLabels[name] = append(m.counterLabels[name], labels)
	return nil
}

func (m *MockMetrics) AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.histogramValues[name] = append(m.histogramValues[name], value)
	m.histogramLabels[name] = append(m.histogramLabels[name], labels)
	return nil
}

func (m *MockMetrics) RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error {
	return nil
}

func (m *MockMetrics) RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error {
	return nil
}

func (m *MockMetrics) RegisterHistogram(name, description, unit string, buckets []float64) error {
	return nil
}

func (m *MockMetrics) GetHistogramValues(name string) []float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.histogramValues[name]
}

func (m *MockMetrics) GetHistogramLabels(name string) []map[string]attribute.Value {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.histogramLabels[name]
}

func (m *MockMetrics) GetCounterValue(name string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counterValues[name]
}

func (m *MockMetrics) GetCounterLabels(name string) []map[string]attribute.Value {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counterLabels[name]
}
