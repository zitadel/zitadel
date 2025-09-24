package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganization_Keys(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		org      *Organization
		index    OrgCacheIndex
		wantKeys []string
	}{
		{
			name:     "index ID returns org ID",
			org:      &Organization{ID: "org1"},
			index:    orgCacheIndexID,
			wantKeys: []string{"org1"},
		},
		{
			name:     "undefined index returns nil",
			org:      &Organization{ID: "org1"},
			index:    orgCacheIndexUndefined,
			wantKeys: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.org.Keys(tc.index)
			assert.ElementsMatch(t, tc.wantKeys, got)
		})
	}
}
