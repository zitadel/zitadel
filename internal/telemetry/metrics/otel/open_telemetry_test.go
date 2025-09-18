package otel

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	colmetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestExporterBehavior(t *testing.T) {
	testCounterName := "my_app_counter"
	testCounterValue := int64(25)

	testCases := []struct {
		name   string
		config map[string]interface{}
		verify func(t *testing.T, m *Metrics, mockData interface{})
		setup  func(t *testing.T) (addr string, mockData interface{})
	}{
		{
			name: "Prometheus Scrape Exporter",
			config: map[string]interface{}{
				"metername": "test-prometheus",
				"registry":  prometheus.NewRegistry(),
			},
			verify: func(t *testing.T, m *Metrics, _ interface{}) {
				server := httptest.NewServer(m.GetExporter())
				defer server.Close()

				resp, err := http.Get(server.URL)
				require.NoError(t, err)
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				bodyStr := string(body)

				assert.Contains(t, bodyStr, testCounterName+"_total 25")
			},
		},
		{
			name:   "OTLP Push Exporter",
			config: make(map[string]interface{}),
			setup: func(t *testing.T) (string, interface{}) {
				metricsCh := make(chan *metricpb.ResourceMetrics, 1)
				addr := startMockCollector(t, metricsCh)
				return addr, metricsCh
			},
			verify: func(t *testing.T, m *Metrics, mockData interface{}) {
				metricsCh := mockData.(chan *metricpb.ResourceMetrics)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				sdkProvider, ok := m.Provider.(*sdk_metric.MeterProvider)
				require.True(t, ok)
				err := sdkProvider.ForceFlush(ctx)
				require.NoError(t, err)

				select {
				case <-metricsCh:
				case <-time.After(2 * time.Second):
					t.Fatal("timed out waiting for metrics from collector")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := tc.config
			var mockData interface{}
			if tc.setup != nil {
				addr, data := tc.setup(t)
				config["endpoint"] = addr
				mockData = data
			}

			m, err := NewMetrics(config)
			require.NoError(t, err)
			metrics, ok := m.(*Metrics)
			require.True(t, ok)

			err = metrics.RegisterCounter(testCounterName, "A test counter")
			require.NoError(t, err)
			err = metrics.AddCount(context.Background(), testCounterName, testCounterValue, nil)
			require.NoError(t, err)

			if tc.verify != nil {
				tc.verify(t, metrics, mockData)
			}
		})
	}
}

type mockCollector struct {
	colmetricpb.UnimplementedMetricsServiceServer
	metricsCh chan<- *metricpb.ResourceMetrics
}

func (c *mockCollector) Export(_ context.Context, req *colmetricpb.ExportMetricsServiceRequest) (*colmetricpb.ExportMetricsServiceResponse, error) {
	for _, rm := range req.ResourceMetrics {
		c.metricsCh <- rm
	}
	return &colmetricpb.ExportMetricsServiceResponse{}, nil
}

func startMockCollector(t *testing.T, metricsCh chan<- *metricpb.ResourceMetrics) string {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	colmetricpb.RegisterMetricsServiceServer(server, &mockCollector{metricsCh: metricsCh})

	go func() {
		if err := server.Serve(lis); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			require.NoError(t, err)
		}
	}()

	t.Cleanup(server.Stop)
	return lis.Addr().String()
}
