//go:build integration

package handlers_test

import (
	"context"
	"os"
	"testing"
	"time"
)

var (
	CTX context.Context
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		CTX = ctx
		return m.Run()
	}())
}
