package admin

import (
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func Test_eventRequestToFilter_ordering(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "instanceID")

	tests := []struct {
		name     string
		req      *admin_pb.ListEventsRequest
		wantDesc bool
	}{
		{
			name:     "default is desc by creation date",
			req:      &admin_pb.ListEventsRequest{},
			wantDesc: true,
		},
		{
			name:     "asc keeps creation date ordering",
			req:      &admin_pb.ListEventsRequest{Asc: true},
			wantDesc: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder, err := eventRequestToFilter(ctx, tt.req)
			if err != nil {
				t.Fatalf("eventRequestToFilter() error = %v", err)
			}
			if !builder.GetOrderByCreationDate() {
				t.Error("expected events to be ordered by creation date")
			}
			if got := builder.GetDesc(); got != tt.wantDesc {
				t.Errorf("GetDesc() = %v, want %v", got, tt.wantDesc)
			}
		})
	}
}

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
