package repository_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
)

func assertMetadata(t *testing.T, expected, gotten []*domain.Metadata) {
	t.Helper()

	assert.Len(t, gotten, len(expected), "metadata length mismatch")
	for _, exp := range expected {
		var actual *domain.Metadata
		gotten = slices.DeleteFunc(gotten, func(m *domain.Metadata) bool {
			if exp.Key != m.Key || exp.InstanceID != m.InstanceID {
				return false
			}
			actual = m
			return true
		})
		require.NotNil(t, actual, "metadata with key %s and instance %s not found", exp.Key, exp.InstanceID)
		assert.Equal(t, exp.Value, actual.Value, "metadata value mismatch for key %s and instance %s", exp.Key, exp.InstanceID)
		assert.NotZero(t, actual.CreatedAt, "metadata created at is zero for key %s and instance %s", exp.Key, exp.InstanceID)
		assert.NotZero(t, actual.UpdatedAt, "metadata updated at is zero for key %s and instance %s", exp.Key, exp.InstanceID)

	}
	assert.Empty(t, gotten, "unmatched metadata found")
}
