package model

import (
	"testing"
)

func TestOnlyLocalhostIsHttp(t *testing.T) {
	type args struct {
		uris []string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "only https uris",
			args: args{
				uris: []string{
					"https://zitadel.ch",
					"https://caos.ch",
				},
			},
			result: true,
		},
		{
			name: "http localhost uris",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://localhost",
				},
			},
			result: true,
		},
		{
			name: "http not localhsot",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://caos.ch",
				},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := onlyLocalhostIsHttp(tt.args.uris)
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
