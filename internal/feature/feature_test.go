package feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	tests := []string{
		"unspecified",
		"login_default_org",
		"trigger_introspection_projections",
		"legacy_introspection",
	}
	for _, want := range tests {
		t.Run(want, func(t *testing.T) {
			feature, err := KeyString(want)
			require.NoError(t, err)
			assert.Equal(t, want, feature.String())
		})
	}
}

func TestLevel(t *testing.T) {
	tests := []string{
		"unspecified",
		"system",
		"instance",
		"org",
		"project",
		"app",
		"user",
	}
	for _, want := range tests {
		t.Run(want, func(t *testing.T) {
			level, err := LevelString(want)
			require.NoError(t, err)
			assert.Equal(t, want, level.String())
		})
	}
}
