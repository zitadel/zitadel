package command

import (
	"context"
	"errors"
	"testing"

	"github.com/caos/zitadel/internal/eventstore"
)

var (
	testSetCommand = func(c *commandBuilder) error {
		c.command = testCommand
		return nil
	}
	testErrSetter = func(c *commandBuilder) error {
		c.err = testErr
		return testErr
	}
	testCommand = func(context.Context) ([]eventstore.Command, error) {
		return nil, nil
	}
	testAgg = &eventstore.Aggregate{Type: "test"}
	testErr = errors.New("test err")
)

func TestNewCommander(t *testing.T) {
	type args struct {
		agg  *eventstore.Aggregate
		opts []validation
	}
	tests := []struct {
		name string
		args args
		want *commandBuilder
	}{
		{
			name: "no aggregate",
			args: args{
				agg:  nil,
				opts: []validation{testSetCommand},
			},
			want: &commandBuilder{
				err:     ErrNoAggregate,
				command: testCommand,
				agg:     nil,
			},
		},
		{
			name: "not executable",
			args: args{
				agg: testAgg,
			},
			want: &commandBuilder{
				agg: testAgg,
				err: ErrNotExecutable,
			},
		},
		{
			name: "with command",
			args: args{
				agg:  testAgg,
				opts: []validation{testSetCommand},
			},
			want: &commandBuilder{
				agg:     testAgg,
				command: testCommand,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCommander(t, tt.want, prepareCommands(tt.args.agg, tt.args.opts...))
		})
	}
}

func Test_commander_Next(t *testing.T) {
	type fields struct {
		err      error
		agg      *eventstore.Aggregate
		previous *commandBuilder
		command  createCommands
	}
	tests := []struct {
		name   string
		fields fields
		args   []validation
		want   *commandBuilder
	}{
		{
			name: "existing error",
			fields: fields{
				err: testErr,
			},
			want: &commandBuilder{
				err: testErr,
			},
		},
		{
			name: "correct",
			fields: fields{
				agg: testAgg,
			},
			args: []validation{
				testSetCommand,
			},
			want: &commandBuilder{
				agg:     testAgg,
				command: testCommand,
				previous: &commandBuilder{
					agg: testAgg,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commandBuilder{
				err:      tt.fields.err,
				agg:      tt.fields.agg,
				previous: tt.fields.previous,
				command:  tt.fields.command,
			}
			assertCommander(t, tt.want, c.Next(tt.args...))
		})
	}
}

func Test_commander_use(t *testing.T) {
	tests := []struct {
		name string
		args []validation
		want *commandBuilder
	}{
		{
			name: "no aggregate",
			args: []validation{
				WithAggregate(nil),
				testSetCommand,
			},
			want: &commandBuilder{
				err:     ErrNoAggregate,
				command: testCommand,
				agg:     nil,
			},
		},
		{
			name: "not executable",
			args: []validation{
				WithAggregate(testAgg),
			},
			want: &commandBuilder{
				agg: testAgg,
				err: ErrNotExecutable,
			},
		},
		{
			name: "with error",
			args: []validation{
				WithAggregate(testAgg),
				testErrSetter,
			},
			want: &commandBuilder{
				agg: testAgg,
				err: testErr,
			},
		},
		{
			name: "with command",
			args: []validation{
				WithAggregate(testAgg),
				testSetCommand,
			},
			want: &commandBuilder{
				agg:     testAgg,
				command: testCommand,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCommander(t, tt.want, (&commandBuilder{}).use(tt.args))
		})
	}
}

func Test_commander_Error(t *testing.T) {
	tests := []struct {
		name      string
		commander *commandBuilder
		err       error
	}{
		{
			name: "no error",
			commander: &commandBuilder{
				previous: &commandBuilder{
					previous: &commandBuilder{},
				},
			},
		},
		{
			name: "error in first",
			commander: &commandBuilder{
				previous: &commandBuilder{
					previous: &commandBuilder{
						err: testErr,
					},
				},
			},
			err: testErr,
		},
		{
			name: "error in last",
			commander: &commandBuilder{
				err: testErr,
				previous: &commandBuilder{
					previous: &commandBuilder{},
				},
			},
			err: testErr,
		},
		{
			name: "error between",
			commander: &commandBuilder{
				previous: &commandBuilder{
					err:      testErr,
					previous: &commandBuilder{},
				},
			},
			err: testErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.commander.Error(); !errors.Is(err, tt.err) {
				t.Errorf("commander.Error() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func assertCommander(t *testing.T, want, got *commandBuilder) {
	t.Helper()
	for want != nil {
		if want.err != got.err {
			t.Errorf("err = %v, want %v", got.err, want.err)
		}
		if want.agg != got.agg {
			t.Errorf("agg = %v, want %v", got.agg, want.agg)
		}
		if (want.previous != nil) != (got.previous != nil) {
			t.Errorf("previous = %v, want %v", got.previous, want.previous)
			return
		}
		if want.previous == nil {
			break
		}
		want = want.previous
		got = got.previous
	}
}
