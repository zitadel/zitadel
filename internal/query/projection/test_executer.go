package projection

import (
	"database/sql"
	"testing"
)

type testExecuter struct {
	expectedStmt string
	gottenStmt   string
	shouldExec   bool

	expectedArgs []interface{}
	gottenArgs   []interface{}
	gotExecuted  bool
}

type anyArg struct{}

func (e *testExecuter) Exec(stmt string, args ...interface{}) (sql.Result, error) {
	e.gottenStmt = stmt
	e.gottenArgs = args
	e.gotExecuted = true
	return nil, nil
}

func (e *testExecuter) Validate(t *testing.T) {
	t.Helper()
	if e.shouldExec != e.gotExecuted {
		t.Error("expected to be executed")
		return
	}
	if len(e.gottenArgs) != len(e.expectedArgs) {
		t.Errorf("wrong arg len expected: %d got: %d", len(e.expectedArgs), len(e.gottenArgs))
	} else {
		for i := 0; i < len(e.expectedArgs); i++ {
			if _, ok := e.expectedArgs[i].(anyArg); ok {
				continue
			}
			if e.expectedArgs[i] != e.gottenArgs[i] {
				t.Errorf("wrong argument at index %d: got: %v want: %v", i, e.gottenArgs[i], e.expectedArgs[i])
			}
		}
	}
	if e.gottenStmt != e.expectedStmt {
		t.Errorf("wrong stmt want:\n%s\ngot:\n%s", e.expectedStmt, e.gottenStmt)
	}
}
