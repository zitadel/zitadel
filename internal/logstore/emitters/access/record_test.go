package access_test

import (
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/internal/logstore/emitters/access"
)

func TestRecord_Normalize(t *testing.T) {
	tests := []struct {
		name   string
		record access.Record
		want   *access.Record
	}{{
		name: "headers with certain keys should be redacted",
		record: access.Record{
			RequestHeaders: map[string][]string{
				"authorization":             {"AValue"},
				"grpcgateway-authorization": {"AValue"},
				"cookie":                    {"AValue"},
				"grpcgateway-cookie":        {"AValue"},
			}, ResponseHeaders: map[string][]string{
				"set-cookie": {"AValue"},
			},
		},
		want: &access.Record{
			RequestHeaders: map[string][]string{
				"authorization":             {"[REDACTED]"},
				"grpcgateway-authorization": {"[REDACTED]"},
				"cookie":                    {"[REDACTED]"},
				"grpcgateway-cookie":        {"[REDACTED]"},
			}, ResponseHeaders: map[string][]string{
				"set-cookie": {"[REDACTED]"},
			},
		},
	}, {
		name: "header keys should be lower cased",
		record: access.Record{
			RequestHeaders:  map[string][]string{"AKey": {"AValue"}},
			ResponseHeaders: map[string][]string{"AKey": {"AValue"}}},
		want: &access.Record{
			RequestHeaders:  map[string][]string{"akey": {"AValue"}},
			ResponseHeaders: map[string][]string{"akey": {"AValue"}}},
	}, {
		name: "an already prune record should stay unchanged",
		record: access.Record{
			RequestURL: "https://my.zitadel.cloud/",
			RequestHeaders: map[string][]string{
				"authorization": {"[REDACTED]"},
			},
			ResponseHeaders: map[string][]string{},
		},
		want: &access.Record{
			RequestURL: "https://my.zitadel.cloud/",
			RequestHeaders: map[string][]string{
				"authorization": {"[REDACTED]"},
			},
			ResponseHeaders: map[string][]string{},
		},
	}, {
		name: "empty record should stay empty",
		record: access.Record{
			RequestHeaders:  map[string][]string{},
			ResponseHeaders: map[string][]string{},
		},
		want: &access.Record{
			RequestHeaders:  map[string][]string{},
			ResponseHeaders: map[string][]string{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.record.Normalize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
