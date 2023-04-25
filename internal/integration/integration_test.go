package integration

import (
	"context"
	"strings"
	"testing"
	"time"
)

const commandLine = `start-from-init --masterkey MasterkeyNeedsToHave32Characters --tlsMode disabled --config ../../e2e/config/localhost/zitadel.yaml --steps ../../e2e/config/localhost/zitadel.yaml`

func TestNewTester(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s := NewTester(ctx, strings.Split(commandLine, " "))
	defer s.Done()
}
