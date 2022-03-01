package command

import (
	"context"
	"errors"
	"testing"

	"github.com/caos/zitadel/internal/eventstore"
)

var (
	testSetCommand = func(c *commander) error {
		c.command = testCommand
		return nil
	}
	testErrSetter = func(c *commander) error {
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
		opts []commanderOption
	}
	tests := []struct {
		name string
		args args
		want *commander
	}{
		{
			name: "no aggregate",
			args: args{
				agg:  nil,
				opts: []commanderOption{testSetCommand},
			},
			want: &commander{
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
			want: &commander{
				agg: testAgg,
				err: ErrNotExecutable,
			},
		},
		{
			name: "with command",
			args: args{
				agg:  testAgg,
				opts: []commanderOption{testSetCommand},
			},
			want: &commander{
				agg:     testAgg,
				command: testCommand,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCommander(t, tt.want, NewCommander(tt.args.agg, tt.args.opts...))
		})
	}
}

func Test_commander_Next(t *testing.T) {
	type fields struct {
		err      error
		agg      *eventstore.Aggregate
		previous *commander
		command  createCommands
	}
	tests := []struct {
		name   string
		fields fields
		args   []commanderOption
		want   *commander
	}{
		{
			name: "existing error",
			fields: fields{
				err: testErr,
			},
			want: &commander{
				err: testErr,
			},
		},
		{
			name: "correct",
			fields: fields{
				agg: testAgg,
			},
			args: []commanderOption{
				testSetCommand,
			},
			want: &commander{
				agg:     testAgg,
				command: testCommand,
				previous: &commander{
					agg: testAgg,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commander{
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
		args []commanderOption
		want *commander
	}{
		{
			name: "no aggregate",
			args: []commanderOption{
				WithAggregate(nil),
				testSetCommand,
			},
			want: &commander{
				err:     ErrNoAggregate,
				command: testCommand,
				agg:     nil,
			},
		},
		{
			name: "not executable",
			args: []commanderOption{
				WithAggregate(testAgg),
			},
			want: &commander{
				agg: testAgg,
				err: ErrNotExecutable,
			},
		},
		{
			name: "with error",
			args: []commanderOption{
				WithAggregate(testAgg),
				testErrSetter,
			},
			want: &commander{
				agg: testAgg,
				err: testErr,
			},
		},
		{
			name: "with command",
			args: []commanderOption{
				WithAggregate(testAgg),
				testSetCommand,
			},
			want: &commander{
				agg:     testAgg,
				command: testCommand,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCommander(t, tt.want, (&commander{}).use(tt.args))
		})
	}
}

func Test_commander_Error(t *testing.T) {
	tests := []struct {
		name      string
		commander *commander
		err       error
	}{
		{
			name: "no error",
			commander: &commander{
				previous: &commander{
					previous: &commander{},
				},
			},
		},
		{
			name: "error in first",
			commander: &commander{
				previous: &commander{
					previous: &commander{
						err: testErr,
					},
				},
			},
			err: testErr,
		},
		{
			name: "error in last",
			commander: &commander{
				err: testErr,
				previous: &commander{
					previous: &commander{},
				},
			},
			err: testErr,
		},
		{
			name: "error between",
			commander: &commander{
				previous: &commander{
					err:      testErr,
					previous: &commander{},
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

func assertCommander(t *testing.T, want, got *commander) {
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
