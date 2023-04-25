//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestNewTester(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s := NewTester(ctx)
	defer s.Done()
}
