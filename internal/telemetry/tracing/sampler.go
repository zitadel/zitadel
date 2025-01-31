package tracing

import (
	"fmt"
	"slices"

	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type spanKindSampler struct {
	sampler sdk_trace.Sampler
	kinds   []trace.SpanKind
}

// ShouldSample implements the [sdk_trace.Sampler] interface.
// It will not sample any spans which do not match the configured span kinds.
// For spans which do match, the decorated sampler is used to make the sampling decision.
func (sk spanKindSampler) ShouldSample(p sdk_trace.SamplingParameters) sdk_trace.SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	if !slices.Contains(sk.kinds, p.Kind) {
		return sdk_trace.SamplingResult{
			Decision:   sdk_trace.Drop,
			Tracestate: psc.TraceState(),
		}
	}
	s := sk.sampler.ShouldSample(p)
	return s
}

func (sk spanKindSampler) Description() string {
	return fmt.Sprintf("SpanKindBased{sampler:%s,kinds:%v}",
		sk.sampler.Description(),
		sk.kinds,
	)
}

// SpanKindBased returns a sampler decorator which behaves differently, based on the kind of the span.
// If the span kind does not match one of the configured kinds, it will not be sampled.
// If the span kind matches, the decorated sampler is used to make sampling decision.
func SpanKindBased(sampler sdk_trace.Sampler, kinds ...trace.SpanKind) sdk_trace.Sampler {
	return spanKindSampler{
		sampler: sampler,
		kinds:   kinds,
	}
}
