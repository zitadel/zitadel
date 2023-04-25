package admin_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

const commandLine = `start-from-init --masterkey MasterkeyNeedsToHave32Characters --tlsMode disabled --config ../../e2e/config/localhost/zitadel.yaml --steps ../../e2e/config/localhost/zitadel.yaml`

var (
	Tester *integration.Tester
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx, strings.Split(commandLine, " "))
		defer Tester.Done()

		return m.Run()
	}())
}
