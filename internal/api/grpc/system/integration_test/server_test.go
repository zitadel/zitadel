//go:build integration

package system_test

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
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour/2)
		defer cancel()

		CTX = ctx
		return m.Run()
	}())
}
