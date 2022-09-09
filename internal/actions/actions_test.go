package actions

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
)

var l logger

type logger struct{}

func (logger) Log(arg string) {
	logging.Info(arg)
}
func (logger) Warn(arg string) {
	logging.Warn(arg)
}
func (logger) Error(arg string) {
	logging.Error(arg)
}

func TestRun(t *testing.T) {
	type args struct {
		api    *API
		script string
		name   string
		opts   []runOpt
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
				opts: []runOpt{},
			},
			wantErr: func(err error) bool { return err == nil },
		},
		{
			name: "throw error",
			args: args{
				api:    nil,
				script: "function testFunc() {throw 'verkackt'}",
				name:   "testFunc",
				opts:   []runOpt{},
			},
			wantErr: func(err error) bool {
				gojaErr := new(goja.Exception)
				if errors.As(err, &gojaErr) {
					return gojaErr.Value().String() == "verkackt"
				}
				return false
			},
		},
		// {
		// 	name: "logger",
		// 	args: args{
		// 		api:    nil,
		// 		script: "",
		// 		name:   "testFunc",
		// 		opts:   []runOpt{
		// 			WithLogger(logger console.Printer)
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "http request",
		// 	args: args{
		// 		api:    nil,
		// 		script: "",
		// 		name:   "testFunc",
		// 		opts:   []runOpt{},
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "interrupt",
			args: args{
				api: nil,
				script: `
function testFunc(){
	let i = 0
	for (;;) {i++}
} 
`,
				name: "testFunc",
				opts: []runOpt{
					WithTimeout(50 * time.Millisecond),
				},
			},
			wantErr: func(err error) bool { return errors.Is(err, ErrHalt) },
		},
		{
			name: "parse response",
			args: args{
				api: nil,
				script: `
let http = require('zitadel/http')
let console = require('zitadel/log')
function testFunc(){
	let res = http.fetch('http://ergast.com/api/f1/2004/1/results.json')
	let parsed = JSON.parse(res.Body)
	// console.warn(res.Body)
	console.warn(JSON.stringify(parsed.MRData))
} 
`,
				name: "testFunc",
				opts: []runOpt{
					WithHTTP(http.DefaultClient),
					WithLogger(l),
				},
			},
			wantErr: func(err error) bool { return errors.Is(err, ErrHalt) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(&Context{}, tt.args.api, tt.args.script, tt.args.name, tt.args.opts...); !tt.wantErr(err) {
				t.Errorf("Run() unexpected error = (%[1]T) %[1]v", err)
			}
		})
	}
}
