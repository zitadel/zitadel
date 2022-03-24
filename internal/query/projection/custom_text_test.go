package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestCustomTextProjection_reduces(t *testing.T) {
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
			name: "org.reduceSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextSetEventType),
					org.AggregateType,
					[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
				), org.CustomTextSetEventMapper),
			},
			reduce: (&CustomTextProjection{}).reduceSet,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       CustomTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO projections.custom_texts (aggregate_id, instance_id, creation_date, change_date, sequence, is_default, template, language, key, text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								false,
								"InitCode",
								"en",
								"Text",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceRemoved",
			reduce: (&CustomTextProjection{}).reduceRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextRemovedEventType),
					org.AggregateType,
					[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
				), org.CustomTextRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       CustomTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts WHERE (aggregate_id = $1) AND (template = $2) AND (key = $3) AND (language = $4)",
							expectedArgs: []interface{}{
								"agg-id",
								"InitCode",
								"Text",
								"en",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceTemplateRemoved",
			reduce: (&CustomTextProjection{}).reduceTemplateRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextTemplateRemovedEventType),
					org.AggregateType,
					[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
				), org.CustomTextTemplateRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       CustomTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts WHERE (aggregate_id = $1) AND (template = $2) AND (language = $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"InitCode",
								"en",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceAdded",
			reduce: (&CustomTextProjection{}).reduceSet,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.CustomTextSetEventType),
					iam.AggregateType,
					[]byte(`{
					"key": "Text",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
				), iam.CustomTextSetEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       CustomTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO projections.custom_texts (aggregate_id, instance_id, creation_date, change_date, sequence, is_default, template, language, key, text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								"InitCode",
								"en",
								"Text",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceRemoved",
			reduce: (&CustomTextProjection{}).reduceRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.CustomTextTemplateRemovedEventType),
					iam.AggregateType,
					[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
				), iam.CustomTextRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       CustomTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts WHERE (aggregate_id = $1) AND (template = $2) AND (key = $3) AND (language = $4)",
							expectedArgs: []interface{}{
								"agg-id",
								"InitCode",
								"Text",
								"en",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceTemplateRemoved",
			reduce: (&CustomTextProjection{}).reduceTemplateRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.CustomTextTemplateRemovedEventType),
					iam.AggregateType,
					[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
				), iam.CustomTextTemplateRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       CustomTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts WHERE (aggregate_id = $1) AND (template = $2) AND (language = $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"InitCode",
								"en",
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
