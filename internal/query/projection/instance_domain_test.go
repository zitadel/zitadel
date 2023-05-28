package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestInstanceDomainProjection_reduces(t *testing.T) {
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
						instance.InstanceDomainAddedEventType,
						instance.AggregateType,
						[]byte(`{"domain": "domain.new", "generated": true}`),
					), instance.DomainAddedEventMapper),
			},
			reduce: (&instanceDomainProjection{}).reduceDomainAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.instance_domains (creation_date, change_date, sequence, domain, instance_id, is_generated, is_primary) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"domain.new",
								"agg-id",
								true,
								false,
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
						instance.InstanceDomainRemovedEventType,
						instance.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), instance.DomainRemovedEventMapper),
			},
			reduce: (&instanceDomainProjection{}).reduceDomainRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_domains WHERE (domain = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(InstanceDomainInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_domains WHERE (instance_id = $1)",
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, InstanceDomainTable, tt.want)
		})
	}
}
