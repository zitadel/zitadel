package eventstore

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func Test_buildPgNotifyQuery(t *testing.T) {
	type args struct {
		events []eventstore.Event
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantArgs  []any
		wantOk    bool
	}{
		{
			name:      "nil events",
			args:      args{nil},
			wantQuery: "",
			wantArgs:  nil,
			wantOk:    false,
		},
		{
			name: "1 event",
			args: args{[]eventstore.Event{
				&event{
					typ:      "foo",
					position: decimal.NewFromInt(1),
				},
			}},
			wantQuery: "SELECT pg_notify($1, $2);",
			wantArgs: []any{
				notificationChannelName,
				"{\"event_type\":\"foo\",\"position\":\"1\"}",
			},
			wantOk: true,
		},
		{
			name: "multiple events",
			args: args{[]eventstore.Event{
				&event{
					typ:      "foo",
					position: decimal.NewFromInt(1),
				},
				&event{
					typ:      "bar",
					position: decimal.NewFromInt(2),
				},
				&event{
					typ:      "spanac",
					position: decimal.NewFromInt(3),
				},
			}},
			wantQuery: "SELECT pg_notify($1, $2), pg_notify($1, $3), pg_notify($1, $4);",
			wantArgs: []any{
				notificationChannelName,
				"{\"event_type\":\"foo\",\"position\":\"1\"}",
				"{\"event_type\":\"bar\",\"position\":\"2\"}",
				"{\"event_type\":\"spanac\",\"position\":\"3\"}",
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotArgs, gotOk := buildPgNotifyQuery(tt.args.events)
			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Equal(t, tt.wantArgs, gotArgs)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}
