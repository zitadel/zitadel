package domain

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
)

func TestUserAgent_GetFingerprintID(t *testing.T) {
	tests := []struct {
		name   string
		fields *UserAgent
		want   string
	}{
		{
			name:   "nil useragent",
			fields: nil,
			want:   "",
		},
		{
			name: "nil fingerprintID",
			fields: &UserAgent{
				FingerprintID: nil,
			},
			want: "",
		},
		{
			name: "value",
			fields: &UserAgent{
				FingerprintID: gu.Ptr("fp"),
			},
			want: "fp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.GetFingerprintID()
			assert.Equal(t, tt.want, got)
		})
	}
}
