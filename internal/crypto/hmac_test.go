package crypto

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHMACValue(t *testing.T) {
	const payload = "123456"

	got := NewHMACValue(payload)
	require.NotNil(t, got)
	assert.Len(t, got.Salt, 16)
	assert.Len(t, got.Hash, sha256.Size)
	assert.True(t, got.Equal(payload), "the returned value must match the payload it was created from")

	other := NewHMACValue(payload)
	assert.NotEqual(t, got.Salt, other.Salt, "the salt must be random")
	assert.NotEqual(t, got.Hash, other.Hash, "an equal payload must not result in an equal hash")
	assert.True(t, other.Equal(payload))
}

func TestNewHMACValue_emptyPayload(t *testing.T) {
	got := NewHMACValue("")
	require.NotNil(t, got)
	assert.Len(t, got.Salt, 16)
	assert.Len(t, got.Hash, sha256.Size)
	assert.True(t, got.Equal(""))
	assert.False(t, got.Equal("123456"))
}

func TestHMACValue_Equal(t *testing.T) {
	const payload = "123456"

	// tampered copies of a valid value.
	wrongHash := NewHMACValue(payload)
	wrongHash.Hash[0]++
	wrongSalt := NewHMACValue(payload)
	wrongSalt.Salt[0]++
	truncatedHash := NewHMACValue(payload)
	truncatedHash.Hash = truncatedHash.Hash[:sha256.Size-1]

	tests := []struct {
		name    string
		value   *HMACValue
		payload string
		want    bool
	}{
		{
			name:    "same payload",
			value:   NewHMACValue(payload),
			payload: payload,
			want:    true,
		},
		{
			name:    "different payload",
			value:   NewHMACValue(payload),
			payload: "654321",
			want:    false,
		},
		{
			name:    "empty payload",
			value:   NewHMACValue(payload),
			payload: "",
			want:    false,
		},
		{
			name:    "payload with an appended character",
			value:   NewHMACValue(payload),
			payload: payload + "0",
			want:    false,
		},
		{
			name:    "modified hash",
			value:   wrongHash,
			payload: payload,
			want:    false,
		},
		{
			name:    "modified salt",
			value:   wrongSalt,
			payload: payload,
			want:    false,
		},
		{
			name:    "truncated hash",
			value:   truncatedHash,
			payload: payload,
			want:    false,
		},
		{
			name:    "zero value",
			value:   new(HMACValue),
			payload: payload,
			want:    false,
		},
		{
			name:    "zero value, empty payload",
			value:   new(HMACValue),
			payload: "",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.value.Equal(tt.payload))
		})
	}
}

// TestHMACValue_json makes sure a value survives a round-trip through an
// event payload, as it is stored and read back as JSON.
func TestHMACValue_json(t *testing.T) {
	const payload = "123456"

	data, err := json.Marshal(NewHMACValue(payload))
	require.NoError(t, err)

	got := new(HMACValue)
	require.NoError(t, json.Unmarshal(data, got))

	assert.Len(t, got.Salt, 16)
	assert.Len(t, got.Hash, sha256.Size)
	assert.True(t, got.Equal(payload), "an unmarshalled value must still match its payload")
	assert.False(t, got.Equal("654321"))
}
