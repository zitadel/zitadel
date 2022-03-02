package org

import (
	"context"
	errs "errors"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestAddDomainCommand(t *testing.T) {
	type args struct {
		a      *org.Aggregate
		domain string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid domain",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "",
			},
			want: want{
				validationErr: errors.ThrowInvalidArgument(nil, "ORG-r3h4J", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "domain",
			},
			want: want{
				commands: []eventstore.Command{
					org.NewDomainAddedEvent(context.Background(), &org.NewAggregate("test", "test").Aggregate, "domain"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidation(t, AddDomainCommand(tt.args.a, tt.args.domain), tt.want)
		})
	}
}

type want struct {
	validationErr error
	createErr     error
	commands      []eventstore.Command
}

func assertValidation(t *testing.T, validation preparation.Validation, want want) {
	t.Helper()

	creates, err := validation()
	if !errs.Is(err, want.validationErr) {
		t.Errorf("wrong validation err = %v, want %v", err, want.validationErr)
		return
	}
	if want.validationErr != nil {
		return
	}
	cmds, err := creates(context.Background(), nil)
	if !errs.Is(err, want.createErr) {
		t.Errorf("wrong create err = %v, want %v", err, want.createErr)
		return
	}
	if want.createErr != nil {
		return
	}

	if len(cmds) != len(want.commands) {
		t.Errorf("wrong length of commands = %d, want %d", len(cmds), len(want.commands))
		return
	}

	for i, cmd := range want.commands {
		if !reflect.DeepEqual(cmd, cmds[i]) {
			t.Errorf("unexpected command: = %v, want %v", cmds[i], cmd)
		}
	}
}
