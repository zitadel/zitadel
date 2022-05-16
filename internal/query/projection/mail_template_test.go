package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
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
				event: getEvent(testEvent(
					repository.EventType(org.MailTemplateAddedEventType),
					org.AggregateType,
					[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
					}`),
				), org.MailTemplateAddedEventMapper),
			},
			reduce: (&mailTemplateProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MailTemplateTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.mail_templates (aggregate_id, creation_date, change_date, sequence, state, is_default, template) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"agg-id",
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
				event: getEvent(testEvent(
					repository.EventType(org.MailTemplateChangedEventType),
					org.AggregateType,
					[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
		}`),
				), org.MailTemplateChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MailTemplateTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.mail_templates SET (change_date, sequence, template) = ($1, $2, $3) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								[]byte("<table></table>"),
								"agg-id",
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
				event: getEvent(testEvent(
					repository.EventType(org.MailTemplateRemovedEventType),
					org.AggregateType,
					nil,
				), org.MailTemplateRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MailTemplateTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.mail_templates WHERE (aggregate_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceAdded",
			reduce: (&mailTemplateProjection{}).reduceAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.MailTemplateAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
					}`),
				), iam.MailTemplateAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       MailTemplateTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.mail_templates (aggregate_id, creation_date, change_date, sequence, state, is_default, template) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name:   "iam.reduceChanged",
			reduce: (&mailTemplateProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.MailTemplateChangedEventType),
					iam.AggregateType,
					[]byte(`{
						"template": "PHRhYmxlPjwvdGFibGU+"
					}`),
				), iam.MailTemplateChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       MailTemplateTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.mail_templates SET (change_date, sequence, template) = ($1, $2, $3) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								[]byte("<table></table>"),
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
			assertReduce(t, got, err, tt.want)
		})
	}
}
