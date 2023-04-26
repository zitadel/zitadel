//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestNewTester(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s := NewTester(ctx)
	defer s.Done()
}
