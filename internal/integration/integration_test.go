//go:build integration

package integration

import (
	"testing"
	"time"
)

func TestNewTester(t *testing.T) {
	ctx, _, cancel := Contexts(time.Hour)
	defer cancel()

	s := NewTester(ctx)
	defer s.Done()
}
