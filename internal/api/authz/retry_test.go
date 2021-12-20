package authz

import (
	"errors"
	"testing"
)

func Test_retry(t *testing.T) {
	type args struct {
		retriable func(*int) func() error
	}
	type want struct {
		executions int
		err        bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 execution",
			args: args{
				retriable: func(execs *int) func() error {
					return func() error {
						if *execs < 1 {
							*execs++
							return errors.New("not 1")
						}
						return nil
					}
				},
			},
			want: want{
				err:        false,
				executions: 1,
			},
		},
		{
			name: "2 execution",
			args: args{
				retriable: func(execs *int) func() error {
					return func() error {
						if *execs < 2 {
							*execs++
							return errors.New("not 2")
						}
						return nil
					}
				},
			},
			want: want{
				err:        false,
				executions: 2,
			},
		},
		{
			name: "too many execution",
			args: args{
				retriable: func(execs *int) func() error {
					return func() error {
						if *execs < 3 {
							*execs++
							return errors.New("not 3")
						}
						return nil
					}
				},
			},
			want: want{
				err:        true,
				executions: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var execs int

			if err := retry(tt.args.retriable(&execs)); (err != nil) != tt.want.err {
				t.Errorf("retry() error = %v, want.err %v", err, tt.want.err)
			}
			if execs != tt.want.executions {
				t.Errorf("retry() executions: want: %d got: %d", tt.want.executions, execs)
			}
		})
	}
}
