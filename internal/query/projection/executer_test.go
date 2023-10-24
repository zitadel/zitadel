package projection

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/errors"
)

type testExecuter struct {
	execIdx    int
	executions []execution
}

type execution struct {
	expectedStmt string
	gottenStmt   string

	expectedArgs []interface{}
	gottenArgs   []interface{}
}

type anyArg struct{}

func (e *testExecuter) Exec(stmt string, args ...interface{}) (sql.Result, error) {
	if stmt == "SAVEPOINT stmt_exec" || stmt == "RELEASE SAVEPOINT stmt_exec" {
		return nil, nil
	}

	if e.execIdx >= len(e.executions) {
		return nil, errors.ThrowInternal(nil, "PROJE-8TNoE", "too many executions")
	}
	e.executions[e.execIdx].gottenArgs = args
	e.executions[e.execIdx].gottenStmt = stmt
	e.execIdx++
	return nil, nil
}

func (e *testExecuter) Validate(t *testing.T) {
	t.Helper()
	if e.execIdx != len(e.executions) {
		t.Errorf("not all expected execs executed. got: %d, want: %d", e.execIdx, len(e.executions))
		return
	}
	for _, execution := range e.executions {
		if len(execution.gottenArgs) != len(execution.expectedArgs) {
			t.Errorf("wrong arg len expected: %d got: %d", len(execution.expectedArgs), len(execution.gottenArgs))
		} else {
			for i := 0; i < len(execution.expectedArgs); i++ {
				if _, ok := execution.expectedArgs[i].(anyArg); ok {
					continue
				}
				assert.Equal(t, execution.expectedArgs[i], execution.gottenArgs[i], "wrong argument at index %d", i)
			}
		}
		if execution.gottenStmt != execution.expectedStmt {
			t.Errorf("wrong stmt want:\n%s\ngot:\n%s", execution.expectedStmt, execution.gottenStmt)
		}
	}
}
