package actions

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/dop251/goja"
// )

// func TestRun(t *testing.T) {
// 	type args struct {
// 		timeout time.Duration
// 		api     *apiParam
// 		script  string
// 		name    string
// 		opts    []Option
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr func(error) bool
// 	}{
// 		{
// 			name: "simple script",
// 			args: args{
// 				api: nil,
// 				script: `
// function testFunc() {
// 	for (i = 0; i < 10; i++) {}
// }`,
// 				name: "testFunc",
// 				opts: []Option{},
// 			},
// 			wantErr: func(err error) bool { return err == nil },
// 		},
// 		{
// 			name: "throw error",
// 			args: args{
// 				api:    nil,
// 				script: "function testFunc() {throw 'verkackt'}",
// 				name:   "testFunc",
// 				opts:   []Option{},
// 			},
// 			wantErr: func(err error) bool {
// 				gojaErr := new(goja.Exception)
// 				if errors.As(err, &gojaErr) {
// 					return gojaErr.Value().String() == "verkackt"
// 				}
// 				return false
// 			},
// 		},
// 		// {
// 		// 	name: "logger",
// 		// 	args: args{
// 		// 		api:    nil,
// 		// 		script: "",
// 		// 		name:   "testFunc",
// 		// 		opts:   []runOpt{
// 		// 			WithLogger(logger console.Printer)
// 		// 		},
// 		// 	},
// 		// 	wantErr: false,
// 		// },
// 		// {
// 		// 	name: "http request",
// 		// 	args: args{
// 		// 		api:    nil,
// 		// 		script: "",
// 		// 		name:   "testFunc",
// 		// 		opts:   []runOpt{},
// 		// 	},
// 		// 	wantErr: false,
// 		// },
// 		// 		{
// 		// 			name: "parse response",
// 		// 			args: args{
// 		// 				api: nil,
// 		// 				script: `
// 		// let http = require('zitadel/http')
// 		// let console = require('zitadel/log')
// 		// function testFunc(){
// 		// 	let res = http.fetch('http://ergast.com/api/f1/2004/1/results.json')
// 		// 	let parsed = JSON.parse(res.Body)
// 		// 	// console.warn(res.Body)
// 		// 	console.warn(JSON.stringify(parsed.MRData))
// 		// }
// 		// `,
// 		// 				name: "testFunc",
// 		// 				opts: []Option{
// 		// 					WithHTTP(http.DefaultClient),
// 		// 					WithLogger(l),
// 		// 				},
// 		// 			},
// 		// 			wantErr: func(err error) bool { return errors.Is(err, ErrHalt) },
// 		// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.args.timeout == 0 {
// 				tt.args.timeout = 10 * time.Second
// 			}
// 			ctx, cancel := context.WithTimeout(context.Background(), tt.args.timeout)
// 			if err := Run(ctx, &Context{}, tt.args.api, tt.args.script, tt.args.name, tt.args.opts...); !tt.wantErr(err) {
// 				t.Errorf("Run() unexpected error = (%[1]T) %[1]v", err)
// 			}
// 			cancel()
// 		})
// 	}
// }
