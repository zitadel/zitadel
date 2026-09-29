package eventstore

import (
	"strings"
	"testing"
)

func TestAddConstraintStmts(t *testing.T) {
	if strings.Contains(addConstraintWithoutOwnersStmt, "owners") {
		t.Fatalf("without-owners insert must not name owners: %s", addConstraintWithoutOwnersStmt)
	}
	if !strings.Contains(addConstraintStmt, "owners") {
		t.Fatalf("tagged insert must name owners: %s", addConstraintStmt)
	}
}
