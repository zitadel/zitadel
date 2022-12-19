package projection

import (
	"context"
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestSearchInstanceDomain_Reduce(t *testing.T) {
	domain := "login.testcorp.ch"
	type fields struct {
		domain string
	}
	type args struct {
		events []eventstore.Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SearchInstanceDomain
	}{
		{
			name: "last instance",
			fields: fields{
				domain: domain,
			},
			want: &SearchInstanceDomain{
				InstanceID: "ok",
			},
			args: args{
				events: []eventstore.Event{
					instance.NewDomainRemovedEvent(
						context.Background(),
						&instance.NewAggregate("1").Aggregate,
						domain,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("1").Aggregate,
						domain,
						false,
					),
					instance.NewDomainRemovedEvent(
						context.Background(),
						&instance.NewAggregate("2").Aggregate,
						domain,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("2").Aggregate,
						domain,
						false,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
						false,
					),
				},
			},
		},
		{
			name: "first instance",
			fields: fields{
				domain: domain,
			},
			want: &SearchInstanceDomain{
				InstanceID: "ok",
			},
			args: args{
				events: []eventstore.Event{
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
						false,
					),
					instance.NewDomainRemovedEvent(
						context.Background(),
						&instance.NewAggregate("1").Aggregate,
						domain,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("1").Aggregate,
						domain,
						false,
					),
					instance.NewDomainRemovedEvent(
						context.Background(),
						&instance.NewAggregate("2").Aggregate,
						domain,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("2").Aggregate,
						domain,
						false,
					),
				},
			},
		},
		{
			name: "added removed added",
			fields: fields{
				domain: domain,
			},
			want: &SearchInstanceDomain{
				InstanceID: "ok",
			},
			args: args{
				events: []eventstore.Event{
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
						false,
					),
					instance.NewDomainRemovedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
						false,
					),
				},
			},
		},
		{
			name: "removed",
			fields: fields{
				domain: domain,
			},
			want: &SearchInstanceDomain{
				InstanceID: "",
			},
			args: args{
				events: []eventstore.Event{
					instance.NewDomainRemovedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
						false,
					),
				},
			},
		},
		{
			name: "no events",
			fields: fields{
				domain: domain,
			},
			want: &SearchInstanceDomain{
				InstanceID: "",
			},
			args: args{},
		},
		{
			name: "livios destroyer",
			fields: fields{
				domain: domain,
			},
			want: &SearchInstanceDomain{
				InstanceID: "",
			},
			args: args{
				events: []eventstore.Event{
					instance.NewInstanceRemovedEvent(
						context.Background(),
						&instance.NewAggregate("1").Aggregate,
						"removed",
						nil,
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("1").Aggregate,
						domain,
						false,
					),
				},
			},
		},
		{
			name: "removed different domain",
			fields: fields{
				domain: domain,
			},
			want: &SearchInstanceDomain{
				InstanceID: "ok",
			},
			args: args{
				events: []eventstore.Event{
					instance.NewDomainRemovedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						"delete.me",
					),
					instance.NewDomainAddedEvent(
						context.Background(),
						&instance.NewAggregate("ok").Aggregate,
						domain,
						false,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domains := NewSearchInstanceDomain(tt.fields.domain)
			domains.Reduce(tt.args.events)

			if tt.want.InstanceID != domains.InstanceID {
				t.Errorf("unexpected instance id. want %q got %q", tt.want.InstanceID, domains.InstanceID)
			}
		})
	}
}
