package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
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
			name: "org reduceSet",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextSetEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), org.CustomTextSetEventMapper),
			},
			reduce: (&customTextProjection{}).reduceSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.custom_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, is_default, template, language, key, text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (instance_id, aggregate_id, template, key, language) DO UPDATE SET (creation_date, change_date, sequence, is_default, text) = (projections.custom_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.is_default, EXCLUDED.text)",
							expectedArgs: []any{
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
			name:   "org reduceRemoved",
			reduce: (&customTextProjection{}).reduceRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts2 WHERE (aggregate_id = $1) AND (template = $2) AND (key = $3) AND (language = $4) AND (instance_id = $5)",
							expectedArgs: []any{
								"agg-id",
								"InitCode",
								"Text",
								"en",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceTemplateRemoved",
			reduce: (&customTextProjection{}).reduceTemplateRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextTemplateRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextTemplateRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts2 WHERE (aggregate_id = $1) AND (template = $2) AND (language = $3) AND (instance_id = $4)",
							expectedArgs: []any{
								"agg-id",
								"InitCode",
								"en",
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
			reduce: reduceInstanceRemovedHelper(CustomTextInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts2 WHERE (instance_id = $1)",
							expectedArgs: []any{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceAdded",
			reduce: (&customTextProjection{}).reduceSet,
			args: args{
				event: getEvent(
					testEvent(
						instance.CustomTextSetEventType,
						instance.AggregateType,
						[]byte(`{
					"key": "Text",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), instance.CustomTextSetEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.custom_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, is_default, template, language, key, text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (instance_id, aggregate_id, template, key, language) DO UPDATE SET (creation_date, change_date, sequence, is_default, text) = (projections.custom_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.is_default, EXCLUDED.text)",
							expectedArgs: []any{
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
			name:   "instance reduceRemoved",
			reduce: (&customTextProjection{}).reduceRemoved,
			args: args{
				event: getEvent(
					testEvent(
						instance.CustomTextTemplateRemovedEventType,
						instance.AggregateType,
						[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
					), instance.CustomTextRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts2 WHERE (aggregate_id = $1) AND (template = $2) AND (key = $3) AND (language = $4) AND (instance_id = $5)",
							expectedArgs: []any{
								"agg-id",
								"InitCode",
								"Text",
								"en",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceTemplateRemoved",
			reduce: (&customTextProjection{}).reduceTemplateRemoved,
			args: args{
				event: getEvent(
					testEvent(
						instance.CustomTextTemplateRemovedEventType,
						instance.AggregateType,
						[]byte(`{
						"key": "Text",
						"language": "en",
						"template": "InitCode"
					}`),
					), instance.CustomTextTemplateRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.custom_texts2 WHERE (aggregate_id = $1) AND (template = $2) AND (language = $3) AND (instance_id = $4)",
							expectedArgs: []any{
								"agg-id",
								"InitCode",
								"en",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&customTextProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.custom_texts2 WHERE (instance_id = $1) AND (aggregate_id = $2)",
							expectedArgs: []any{
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}
			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, CustomTextTable, tt.want)
		})
	}
}
