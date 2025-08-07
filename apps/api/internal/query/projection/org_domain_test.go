package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestOrgDomainProjection_reduces(t *testing.T) {
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
						org.OrgDomainAddedEventType,
						org.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), org.DomainAddedEventMapper),
			},
			reduce: (&orgDomainProjection{}).reduceDomainAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.org_domains2 (creation_date, change_date, sequence, domain, org_id, instance_id, is_verified, is_primary, validation_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"domain.new",
								"agg-id",
								"instance-id",
								false,
								false,
								domain.OrgDomainValidationTypeUnspecified,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceDomainVerificationAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgDomainVerificationAddedEventType,
						org.AggregateType,
						[]byte(`{"domain": "domain.new", "validationType": 2}`),
					), org.DomainVerificationAddedEventMapper),
			},
			reduce: (&orgDomainProjection{}).reduceDomainVerificationAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.org_domains2 SET (change_date, sequence, validation_type) = ($1, $2, $3) WHERE (domain = $4) AND (org_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.OrgDomainValidationTypeDNS,
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
			name: "reduceDomainVerified",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgDomainVerifiedEventType,
						org.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), org.DomainVerifiedEventMapper),
			},
			reduce: (&orgDomainProjection{}).reduceDomainVerified,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.org_domains2 SET (change_date, sequence, is_verified) = ($1, $2, $3) WHERE (domain = $4) AND (org_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
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
			name: "reducePrimaryDomainSet",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgDomainPrimarySetEventType,
						org.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), org.DomainPrimarySetEventMapper),
			},
			reduce: (&orgDomainProjection{}).reducePrimaryDomainSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.org_domains2 SET (change_date, sequence, is_primary) = ($1, $2, $3) WHERE (org_id = $4) AND (is_primary = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								false,
								"agg-id",
								true,
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.org_domains2 SET (change_date, sequence, is_primary) = ($1, $2, $3) WHERE (domain = $4) AND (org_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
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
			name: "reduceDomainRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgDomainRemovedEventType,
						org.AggregateType,
						[]byte(`{"domain": "domain.new"}`),
					), org.DomainRemovedEventMapper),
			},
			reduce: (&orgDomainProjection{}).reduceDomainRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.org_domains2 WHERE (domain = $1) AND (org_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
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
			name:   "org.reduceOwnerRemoved",
			reduce: (&orgDomainProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.org_domains2 WHERE (instance_id = $1) AND (org_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
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
			reduce: reduceInstanceRemovedHelper(OrgDomainInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.org_domains2 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, OrgDomainTable, tt.want)
		})
	}
}
