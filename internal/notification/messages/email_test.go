package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmail_GetContent_Cc(t *testing.T) {
	tests := []struct {
		name     string
		cc       []string
		wantCc   bool
		wantLine string
	}{
		{
			name:   "no cc",
			cc:     nil,
			wantCc: false,
		},
		{
			name:   "empty cc",
			cc:     []string{},
			wantCc: false,
		},
		{
			name:     "single cc",
			cc:       []string{"cc@example.com"},
			wantCc:   true,
			wantLine: "Cc: cc@example.com\r\n",
		},
		{
			name:     "multiple cc",
			cc:       []string{"cc1@example.com", "cc2@example.com"},
			wantCc:   true,
			wantLine: "Cc: cc1@example.com, cc2@example.com\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Email{
				Recipients:  []string{"to@example.com"},
				CC:          tt.cc,
				SenderEmail: "from@example.com",
				Subject:     "subject",
				Content:     "content",
			}
			got, err := msg.GetContent()
			require.NoError(t, err)
			if !tt.wantCc {
				assert.NotContains(t, got, "Cc:")
				return
			}
			assert.Contains(t, got, tt.wantLine)
		})
	}
}
