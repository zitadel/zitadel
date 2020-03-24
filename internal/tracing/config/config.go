package config

import (
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/caos/zitadel/internal/tracing"
	tracing_g "github.com/caos/zitadel/internal/tracing/google"
	tracing_log "github.com/caos/zitadel/internal/tracing/log"
)

type TracingConfig struct {
	Type   string
	Config tracing.Config
}

var tracer = map[string]func() tracing.Config{
	"google": func() tracing.Config { return &tracing_g.Config{} },
	"log":    func() tracing.Config { return &tracing_log.Config{} },
}

func (c *TracingConfig) UnmarshalJSON(data []byte) error {
	var rc struct {
		Type   string
		Config json.RawMessage
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return status.Errorf(codes.Internal, "%v parse config: %v", "TRACE-vmjS", err)
	}

	c.Type = rc.Type

	var err error
	c.Config, err = newTracingConfig(c.Type, rc.Config)
	if err != nil {
		return status.Errorf(codes.Internal, "%v parse config: %v", "TRACE-Ws9E", err)
	}

	return c.Config.NewTracer()
}

func newTracingConfig(tracerType string, configData []byte) (tracing.Config, error) {
	t, ok := tracer[tracerType]
	if !ok {
		return nil, status.Errorf(codes.Internal, "%v No config: %v", "TRACE-HMEJ", tracerType)
	}

	tracingConfig := t()
	if len(configData) == 0 {
		return tracingConfig, nil
	}

	if err := json.Unmarshal(configData, tracingConfig); err != nil {
		return nil, status.Errorf(codes.Internal, "%v Could not read conifg: %v", "TRACE-1tSS", err)
	}

	return tracingConfig, nil
}
