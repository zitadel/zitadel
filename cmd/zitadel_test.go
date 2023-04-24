package cmd

import (
	"strings"
	"testing"
)

const commandLine = `start-from-init --masterkey MasterkeyNeedsToHave32Characters --tlsMode disabled --config ../e2e/config/localhost/zitadel.yaml --steps ../e2e/config/localhost/zitadel.yaml`

func TestNewTestServer(t *testing.T) {
	s := NewTestServer(strings.Split(commandLine, " "))
	defer s.Done()
}
