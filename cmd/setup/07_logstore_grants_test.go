package setup

import (
	"strings"
	"testing"
)

// The logstore schema statement runs as the service user itself, so it must
// not contain an unconditional GRANT statement: managed PostgreSQL services
// (e.g. Fly.io Managed Postgres) reject GRANT entirely.
func Test_logstoreStmtHasNoUnconditionalGrant(t *testing.T) {
	for _, line := range strings.Split(createLogstoreSchema07, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "GRANT") {
			t.Errorf("07/logstore.sql contains an unconditional GRANT statement: %q", strings.TrimSpace(line))
		}
	}
	if !strings.Contains(createLogstoreSchema07, "current_user <>") {
		t.Error("07/logstore.sql lost its grant guard; grants must be skipped only when they are self-grants")
	}
}
