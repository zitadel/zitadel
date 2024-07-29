package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestInstanceAllowedDomainProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reduceDomainAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.AllowedDomainAddedEventType,
						instance.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), eventstore.GenericEventMapper[instance.AllowedDomainAddedEvent]),
			},
			reduce: (&instanceAllowedDomainProjection{}).reduceDomainAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.instance_allowed_domains (creation_date, change_date, sequence, domain, instance_id) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"domain.new",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceDomainRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.AllowedDomainRemovedEventType,
						instance.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), eventstore.GenericEventMapper[instance.AllowedDomainRemovedEvent]),
			},
			reduce: (&instanceAllowedDomainProjection{}).reduceDomainRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_allowed_domains WHERE (domain = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"domain.new",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(InstanceAllowedDomainInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_allowed_domains WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, InstanceAllowedDomainTable, tt.want)
		})
	}
}
