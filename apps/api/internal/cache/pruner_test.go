package cache

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

type testPruner struct {
	called chan struct{}
}

func (p *testPruner) Prune(context.Context) error {
	p.called <- struct{}{}
	return nil
}

func TestAutoPruneConfig_startAutoPrune(t *testing.T) {
	c := AutoPruneConfig{
		Interval: time.Second,
		Timeout:  time.Millisecond,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pruner := testPruner{
		called: make(chan struct{}),
	}
	clock := clockwork.NewFakeClock()
	close := c.startAutoPrune(ctx, &pruner, PurposeAuthzInstance, clock)
	defer close()
	clock.Advance(time.Second)

	select {
	case _, ok := <-pruner.called:
		assert.True(t, ok)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}
