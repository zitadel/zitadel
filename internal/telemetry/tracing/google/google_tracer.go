package google

import (
	"os"
	"strings"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/telemetry/tracing/otel"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	ProjectID    string
	MetricPrefix string
	Fraction     float64
}

type Tracer struct {
	otel.Tracer
}

func (c *Config) NewTracer() error {
	if !envIsSet() {
		return errors.ThrowInvalidArgument(nil, "GOOGL-sdh3a", "env not properly set, GOOGLE_APPLICATION_CREDENTIALS is misconfigured or missing")
	}

	sampler := sdk_trace.ParentBased(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := texporter.NewExporter(texporter.WithProjectID(c.ProjectID))
	if err != nil {
		return err
	}

	tracing.T = &Tracer{Tracer: *(otel.NewTracer(c.MetricPrefix, sampler, exporter))}

	return nil
}

func envIsSet() bool {
	gAuthCred := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	return strings.Contains(gAuthCred, ".json")
}
