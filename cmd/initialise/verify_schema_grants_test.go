package initialise

import (
	"fmt"
	"strings"
	"testing"
)

// The schema bootstrap statements run as the service user itself in every
// flow, so they must not contain unconditional GRANT statements: managed
// PostgreSQL services (e.g. Fly.io Managed Postgres) reject GRANT entirely.
func Test_schemaStmtsHaveNoUnconditionalGrants(t *testing.T) {
	if err := ReadStmts(); err != nil {
		t.Fatalf("unable to read stmts: %v", err)
	}

	stmts := map[string]string{
		"04_eventstore":  createEventstoreStmt,
		"05_projections": createProjectionsStmt,
		"06_system":      createSystemStmt,
	}
	for name, stmt := range stmts {
		for _, line := range strings.Split(stmt, "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "GRANT") {
				t.Errorf("%s contains an unconditional GRANT statement: %q", name, strings.TrimSpace(line))
			}
		}
		formatted := fmt.Sprintf(stmt, "service_user")
		if strings.Contains(formatted, "%!") || !strings.Contains(formatted, "TO %I") {
			t.Errorf("%s has invalid formatting after username substitution: %s", name, formatted)
		}
		if !strings.Contains(stmt, "current_user <>") {
			t.Errorf("%s lost its grant guard; grants must be skipped only when they are self-grants", name)
		}
	}
}
