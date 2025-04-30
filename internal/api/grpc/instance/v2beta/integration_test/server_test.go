package instance_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	CTXWithSysAuthZ context.Context
)

func TestMain(m *testing.M) {
	os.Exit(func() int {

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		CTXWithSysAuthZ = integration.WithSystemAuthorization(ctx)

		return m.Run()
	}())
}
