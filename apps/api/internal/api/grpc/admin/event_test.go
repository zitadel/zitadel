package admin

import (
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func Test_aggregateTypesFromEventTypes(t *testing.T) {
	type args struct {
		eventTypes []eventstore.EventType
	}
	tests := []struct {
		name string
		args args
		want []eventstore.AggregateType
	}{
		{
			name: "no event types",
			args: args{
				eventTypes: []eventstore.EventType{},
			},
			want: []eventstore.AggregateType{},
		},
		{
			name: "only by prefix",
			args: args{
				eventTypes: []eventstore.EventType{user.MachineAddedEventType, org.OrgAddedEventType},
			},
			want: []eventstore.AggregateType{user.AggregateType, org.AggregateType},
		},
		{
			name: "with special",
			args: args{
				eventTypes: []eventstore.EventType{deviceauth.ApprovedEventType, org.OrgAddedEventType},
			},
			want: []eventstore.AggregateType{deviceauth.AggregateType, org.AggregateType},
		},
		{
			name: "duplicates",
			args: args{
				eventTypes: []eventstore.EventType{org.OrgAddedEventType, org.OrgChangedEventType},
			},
			want: []eventstore.AggregateType{org.AggregateType, org.AggregateType},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aggregateTypesFromEventTypes(tt.args.eventTypes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("aggregateTypesFromEventTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}
