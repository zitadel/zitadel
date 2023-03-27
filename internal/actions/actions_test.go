package actions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/zitadel/internal/logstore"
)

func TestRun(t *testing.T) {
	SetLogstoreService(logstore.New(nil, nil, nil))
	type args struct {
		timeout time.Duration
		api     apiFields
		ctx     contextFields
		script  string
		name    string
		opts    []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr func(error) bool
	}{
		{
			name: "simple script",
			args: args{
				api: nil,
				script: `
function testFunc() {
	for (i = 0; i < 10; i++) {}
}`,
				name: "testFunc",
				opts: []Option{},
			},
			wantErr: func(err error) bool { return err == nil },
		},
		{
			name: "throw error",
			args: args{
				api:    nil,
				script: "function testFunc() {throw 'some error'}",
				name:   "testFunc",
				opts:   []Option{},
			},
			wantErr: func(err error) bool {
				gojaErr := new(goja.Exception)
				if errors.As(err, &gojaErr) {
					return gojaErr.Value().String() == "some error"
				}
				return false
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.timeout == 0 {
				tt.args.timeout = 10 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.Background(), tt.args.timeout)
			if err := Run(ctx, tt.args.ctx, tt.args.api, tt.args.script, tt.args.name, tt.args.opts...); !tt.wantErr(err) {
				t.Errorf("Run() unexpected error = (%[1]T) %[1]v", err)
			}
			cancel()
		})
	}
}
