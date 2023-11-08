package preparation

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
)

var errTest = errors.New("test")

func Test_validate(t *testing.T) {
	type args struct {
		validations []Validation
	}
	type want struct {
		len int
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no validations",
			args: args{},
			want: want{
				err: ErrNotExecutable,
			},
		},
		{
			name: "error in validation",
			args: args{
				validations: []Validation{
					func() (CreateCommands, error) {
						return nil, errTest
					},
				},
			},
			want: want{
				err: errTest,
			},
		},
		{
			name: "correct",
			args: args{
				validations: []Validation{
					func() (CreateCommands, error) {
						return func(_ context.Context, _ FilterToQueryReducer) ([]eventstore.Command, error) {
							return nil, nil
						}, nil
					},
				},
			},
			want: want{
				len: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate(tt.args.validations)
			if !errors.Is(err, tt.want.err) {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if len(got) != tt.want.len {
				t.Errorf("validate() len = %v, want %v", len(got), tt.want.len)
			}
		})
	}
}

func Test_create(t *testing.T) {
	type args struct {
		commanders []CreateCommands
	}
	type want struct {
		err error
		len int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "error in command",
			want: want{
				err: errTest,
			},
			args: args{
				commanders: []CreateCommands{
					func(_ context.Context, _ FilterToQueryReducer) ([]eventstore.Command, error) {
						return nil, errTest
					},
				},
			},
		},
		{
			name: "no commands",
			want: want{},
			args: args{
				commanders: []CreateCommands{
					func(_ context.Context, _ FilterToQueryReducer) ([]eventstore.Command, error) {
						return nil, nil
					},
				},
			},
		},
		{
			name: "multiple commands",
			want: want{
				len: 3,
			},
			args: args{
				commanders: []CreateCommands{
					func(_ context.Context, _ FilterToQueryReducer) ([]eventstore.Command, error) {
						return []eventstore.Command{new(testCommand), new(testCommand)}, nil
					},
					func(_ context.Context, _ FilterToQueryReducer) ([]eventstore.Command, error) {
						return []eventstore.Command{new(testCommand)}, nil
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmds, err := create(context.Background(), nil, tt.args.commanders)
			if !errors.Is(err, tt.want.err) {
				t.Errorf("create() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if len(gotCmds) != tt.want.len {
				t.Errorf("create() len = %d, want %d", len(gotCmds), tt.want.len)
			}
		})
	}
}

func Test_transactionFilter(t *testing.T) {
	type args struct {
		filter   FilterToQueryReducer
		commands []eventstore.Command
	}
	tests := []struct {
		name string
		args args
		want FilterToQueryReducer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transactionFilter(tt.args.filter, tt.args.commands); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transactionFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testCommand struct {
	eventstore.BaseEvent
}

func (c *testCommand) Payload() interface{} {
	return nil
}

func (c *testCommand) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
