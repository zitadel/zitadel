package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExporterType_isNone(t *testing.T) {
	tests := []struct {
		name string
		e    ExporterType
		want bool
	}{
		{
			name: "unspecified is none",
			e:    ExporterTypeUnspecified,
			want: true,
		},
		{
			name: "none is none",
			e:    ExporterTypeNone,
			want: true,
		},
		{
			name: "stdout is not none",
			e:    ExporterTypeStdOut,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.isNone()
			assert.Equal(t, tt.want, got)
		})
	}
}
