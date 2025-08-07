package schemas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteOnlyString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		s    WriteOnlyString
	}{
		{
			name: "always returns null",
			s:    "foo bar",
		},
		{
			name: "empty string returns null",
			s:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(&tt.s)
			assert.NoError(t, err)
			assert.Equal(t, "null", string(got))
		})
	}
}

func TestWriteOnlyString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    WriteOnlyString
		wantErr bool
	}{
		{
			name:    "string",
			input:   []byte(`"fooBar"`),
			want:    "fooBar",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   []byte(`""`),
			want:    "",
			wantErr: false,
		},
		{
			name:    "bad format",
			input:   []byte(`"bad "format"`),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got WriteOnlyString
			err := json.Unmarshal(tt.input, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
