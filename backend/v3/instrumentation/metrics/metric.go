package metrics

import (
	"context"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/cmd/build"
)

const (
	ActiveSessionCounter            = "zitadel.active_session_counter"
	ActiveSessionCounterDescription = "Active session counter"
	SpoolerDivCounter               = "zitadel.spooler_div_milliseconds"
	SpoolerDivCounterDescription    = "Spooler div from last successful run to now in milliseconds"
	Database                        = "database"
	ViewName                        = "view_name"
)

type Metrics interface {
	GetExporter() http.Handler
	RegisterCounter(name, description string) error
	AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error
	AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error
	RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error
	RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error
	RegisterHistogram(name, description, unit string, buckets []float64) error
}

const pkgName = "github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"

var M = sync.OnceValue(func() Metrics {
	return instrumentation.NewMeter(
		pkgName,
		metric.WithInstrumentationVersion(build.Version()),
	)
})

func GetExporter() http.Handler {
	return M().GetExporter()
}

func RegisterCounter(name, description string) error {
	return M().RegisterCounter(name, description)
}

func AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error {
	return M().AddCount(ctx, name, value, labels)
}

func AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error {
	return M().AddHistogramMeasurement(ctx, name, value, labels)
}

func RegisterHistogram(name, description, unit string, buckets []float64) error {
	return M().RegisterHistogram(name, description, unit, buckets)
}

func RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error {
	return M().RegisterUpDownSumObserver(name, description, callbackFunc)
}

func RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error {
	return M().RegisterValueObserver(name, description, callbackFunc)
}
