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

func TestMailTemplateProjection_reduces(t *testing.T) {
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
			name: "org.reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.MailTemplateAddedEventType,
						org.AggregateType,
						[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
					}`),
					), org.MailTemplateAddedEventMapper),
			},
			reduce: (&mailTemplateProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.mail_templates2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, is_default, template) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								false,
								[]byte("<table></table>"),
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceChanged",
			reduce: (&mailTemplateProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						org.MailTemplateChangedEventType,
						org.AggregateType,
						[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
		}`),
					), org.MailTemplateChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.mail_templates2 SET (change_date, sequence, template) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								[]byte("<table></table>"),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceRemoved",
			reduce: (&mailTemplateProjection{}).reduceRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.MailTemplateRemovedEventType,
						org.AggregateType,
						nil,
					), org.MailTemplateRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.mail_templates2 WHERE (aggregate_id = $1) AND (instance_id = $2)",
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
			name: "instance.reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(MailTemplateInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.mail_templates2 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance.reduceAdded",
			reduce: (&mailTemplateProjection{}).reduceAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.MailTemplateAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
					}`),
					), instance.MailTemplateAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.mail_templates2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, is_default, template) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								true,
								[]byte("<table></table>"),
							},
						},
					},
				},
			},
		},
		{
			name:   "instance.reduceChanged",
			reduce: (&mailTemplateProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						instance.MailTemplateChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
					}`),
					), instance.MailTemplateChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.mail_templates2 SET (change_date, sequence, template) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								[]byte("<table></table>"),
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
			reduce: (&mailTemplateProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.mail_templates2 WHERE (instance_id = $1) AND (aggregate_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
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
			assertReduce(t, got, err, MailTemplateTable, tt.want)
		})
	}
}
