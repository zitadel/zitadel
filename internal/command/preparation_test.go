// this is a helper file for tests
package command

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/eventstore"
)

//Want represents the expected values for each step
type Want struct {
	ValidationErr error
	CreateErr     error
	Commands      []eventstore.Command
}

//AssertValidation checks if the validation works as inteded
func AssertValidation(t *testing.T, validation preparation.Validation, filter preparation.FilterToQueryReducer, want Want) {
	t.Helper()

	creates, err := validation()
	if !errors.Is(err, want.ValidationErr) {
		t.Errorf("wrong validation err = %v, want %v", err, want.ValidationErr)
		return
	}
	if err != nil {
		return
	}
	cmds, err := creates(context.Background(), filter)
	if !errors.Is(err, want.CreateErr) {
		t.Errorf("wrong create err = %v, want %v", err, want.CreateErr)
		return
	}
	if err != nil {
		return
	}

	if len(cmds) != len(want.Commands) {
		t.Errorf("wrong length of commands = %v, want %v", eventTypes(cmds), eventTypes(want.Commands))
		return
	}

	for i, cmd := range want.Commands {
		if !reflect.DeepEqual(cmd, cmds[i]) {
			t.Errorf("unexpected command: = %v, want %v", cmds[i], cmd)
		}
	}
}

func eventTypes(cmds []eventstore.Command) []eventstore.EventType {
	types := make([]eventstore.EventType, len(cmds))
	for i, cmd := range cmds {
		types[i] = cmd.Type()
	}
	return types
}

type MultiFilter struct {
	count   int
	filters []preparation.FilterToQueryReducer
}

func NewMultiFilter() *MultiFilter {
	return new(MultiFilter)
}

func (mf *MultiFilter) Append(filter preparation.FilterToQueryReducer) *MultiFilter {
	mf.filters = append(mf.filters, filter)
	return mf
}

func (mf *MultiFilter) Filter() preparation.FilterToQueryReducer {
	return func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
		mf.count++
		return mf.filters[mf.count-1](ctx, queryFactory)
	}
}
