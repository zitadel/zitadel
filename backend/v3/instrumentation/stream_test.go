package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	EnableStreams(StreamRuntime)
	assert.True(t, IsStreamEnabled(StreamRuntime))
	assert.False(t, IsStreamEnabled(StreamAction))
}
