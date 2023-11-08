package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestOrgProjection_reduces(t *testing.T) {
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
			name: "reducePrimaryDomainSet",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgDomainPrimarySetEventType,
						org.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), org.DomainPrimarySetEventMapper),
			},
			reduce: (&orgProjection{}).reducePrimaryDomainSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.orgs1 SET (change_date, sequence, primary_domain) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"domain.new",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOrgReactivated",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgReactivatedEventType,
						org.AggregateType,
						nil,
					), org.OrgReactivatedEventMapper),
			},
			reduce: (&orgProjection{}).reduceOrgReactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.orgs1 SET (change_date, sequence, org_state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.OrgStateActive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOrgDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgDeactivatedEventType,
						org.AggregateType,
						nil,
					), org.OrgDeactivatedEventMapper),
			},
			reduce: (&orgProjection{}).reduceOrgDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.orgs1 SET (change_date, sequence, org_state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.OrgStateInactive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOrgChanged",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgChangedEventType,
						org.AggregateType,
						[]byte(`{"name": "new name"}`),
					), org.OrgChangedEventMapper),
			},
			reduce: (&orgProjection{}).reduceOrgChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.orgs1 SET (change_date, sequence, name) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"new name",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOrgChanged no changes",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgChangedEventType,
						org.AggregateType,
						[]byte(`{}`),
					), org.OrgChangedEventMapper),
			},
			reduce: (&orgProjection{}).reduceOrgChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer:      &testExecuter{},
			},
		},
		{
			name: "reduceOrgAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgAddedEventType,
						org.AggregateType,
						[]byte(`{"name": "name"}`),
					), org.OrgAddedEventMapper),
			},
			reduce: (&orgProjection{}).reduceOrgAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.orgs1 (id, creation_date, change_date, resource_owner, instance_id, sequence, name, org_state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"name",
								domain.OrgStateActive,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOrgRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&orgProjection{}).reduceOrgRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.orgs1 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
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
			reduce: reduceInstanceRemovedHelper(OrgColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.orgs1 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, OrgProjectionTable, tt.want)
		})
	}
}
