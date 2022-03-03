// this is a helper file for tests
package preparation

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/eventstore"
)

//Want represents the expected values for each step
type Want struct {
	ValidationErr error
	CreateErr     error
	Commands      []eventstore.Command
}

//AssertValidation checks if the validation works as inteded
func AssertValidation(t *testing.T, validation Validation, want Want) {
	t.Helper()

	creates, err := validation()
	if !errors.Is(err, want.ValidationErr) {
		t.Errorf("wrong validation err = %v, want %v", err, want.ValidationErr)
		return
	}
	if err != nil {
		return
	}
	cmds, err := creates(context.Background(), nil)
	if !errors.Is(err, want.CreateErr) {
		t.Errorf("wrong create err = %v, want %v", err, want.CreateErr)
		return
	}
	if err != nil {
		return
	}

	if len(cmds) != len(want.Commands) {
		t.Errorf("wrong length of commands = %d, want %d", len(cmds), len(want.Commands))
		return
	}

	for i, cmd := range want.Commands {
		if !reflect.DeepEqual(cmd, cmds[i]) {
			t.Errorf("unexpected command: = %v, want %v", cmds[i], cmd)
		}
	}
}
