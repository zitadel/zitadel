package cache

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testPruner struct {
	called atomic.Int64
}

func (p *testPruner) Prune(context.Context) error {
	p.called.Add(1)
	return nil
}

func TestAutoPruneConfig_StartAutoPrune(t *testing.T) {
	c := AutoPruneConfig{
		Interval: 300 * time.Millisecond,
		Timeout:  time.Millisecond,
	}
	var pruner testPruner
	close := c.StartAutoPrune(context.Background(), &pruner)
	defer close()

	// 3 runs take 900 milliseconds.
	time.Sleep(time.Second)
	assert.Equal(t, int64(3), pruner.called.Load())
}
