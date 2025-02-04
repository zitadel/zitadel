package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestSpanKindBased(t *testing.T) {
	type args struct {
		sampler sdk_trace.Sampler
		kinds   []trace.SpanKind
	}
	type want struct {
		description string
		sampled     int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"never sample, no sample",
			args{
				sampler: sdk_trace.NeverSample(),
				kinds:   []trace.SpanKind{trace.SpanKindServer},
			},
			want{
				description: "SpanKindBased{sampler:AlwaysOffSampler,kinds:[server]}",
				sampled:     0,
			},
		},
		{
			"always sample, no kind, no sample",
			args{
				sampler: sdk_trace.AlwaysSample(),
				kinds:   nil,
			},
			want{
				description: "SpanKindBased{sampler:AlwaysOnSampler,kinds:[]}",
				sampled:     0,
			},
		},
		{
			"always sample, 2 kinds, 2 samples",
			args{
				sampler: sdk_trace.AlwaysSample(),
				kinds:   []trace.SpanKind{trace.SpanKindServer, trace.SpanKindClient},
			},
			want{
				description: "SpanKindBased{sampler:AlwaysOnSampler,kinds:[server client]}",
				sampled:     2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := SpanKindBased(tt.args.sampler, tt.args.kinds...)
			assert.Equal(t, tt.want.description, sampler.Description())

			p := sdk_trace.NewTracerProvider(sdk_trace.WithSampler(sampler))
			tr := p.Tracer("test")

			var sampled int
			for i := trace.SpanKindUnspecified; i <= trace.SpanKindConsumer; i++ {
				ctx := context.Background()
				_, span := tr.Start(ctx, "test", trace.WithSpanKind(i))
				if span.SpanContext().IsSampled() {
					sampled++
				}
			}

			assert.Equal(t, tt.want.sampled, sampled)
		})
	}
}
