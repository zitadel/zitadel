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

func TestMessageTextProjection_reduces(t *testing.T) {
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
			name: "org reduceAdded Title",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextSetEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Title",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), org.CustomTextSetEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, title) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, title) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.title)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAdded PreHeader",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextSetEventType,
						org.AggregateType,
						[]byte(`{
						"key": "PreHeader",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), org.CustomTextSetEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, pre_header) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, pre_header) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.pre_header)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAdded Subject",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextSetEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Subject",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), org.CustomTextSetEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, subject) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, subject) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.subject)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAdded Greeting",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextSetEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Greeting",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), org.CustomTextSetEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, greeting) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, greeting) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.greeting)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAdded Text",
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
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, text) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.text)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAdded ButtonText",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextSetEventType,
						org.AggregateType,
						[]byte(`{
						"key": "ButtonText",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), org.CustomTextSetEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, button_text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, button_text) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.button_text)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAdded Footer",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextSetEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Footer",
						"language": "en",
						"template": "InitCode",
						"text": "Test"
					}`),
					), org.CustomTextSetEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, footer_text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, footer_text) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.footer_text)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceRemoved Title",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Title",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, title) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			reduce: reduceInstanceRemovedHelper(MessageTextInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.message_texts2 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceRemoved PreHeader",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "PreHeader",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, pre_header) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org reduceRemoved Subject",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Subject",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, subject) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org reduceRemoved Greeting",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Greeting",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, greeting) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org reduceRemoved Text",
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
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, text) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org reduceRemoved ButtonText",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "ButtonText",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, button_text) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org reduceRemoved Footer",
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Footer",
						"language": "en",
						"template": "InitCode"
					}`),
					), org.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, footer_text) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name:   "org reduceRemoved",
			reduce: (&messageTextProjection{}).reduceTemplateRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.CustomTextTemplateRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"key": "Title", 
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
							expectedStmt: "DELETE FROM projections.message_texts2 WHERE (aggregate_id = $1) AND (type = $2) AND (language = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
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
			name:   "instance reduceAdded",
			reduce: (&messageTextProjection{}).reduceAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.CustomTextSetEventType,
						instance.AggregateType,
						[]byte(`{
						"key": "Title",
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
							expectedStmt: "INSERT INTO projections.message_texts2 (aggregate_id, instance_id, creation_date, change_date, sequence, state, type, language, title) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (instance_id, aggregate_id, type, language) DO UPDATE SET (creation_date, change_date, sequence, state, title) = (projections.message_texts2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.title)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.PolicyStateActive,
								"InitCode",
								"en",
								"Test",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceRemoved Title",
			args: args{
				event: getEvent(
					testEvent(
						instance.CustomTextRemovedEventType,
						instance.AggregateType,
						[]byte(`{
						"key": "Title",
						"language": "en",
						"template": "InitCode"
					}`),
					), instance.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.message_texts2 SET (change_date, sequence, title) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			reduce: (&messageTextProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.message_texts2 WHERE (instance_id = $1) AND (aggregate_id = $2)",
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, MessageTextTable, tt.want)
		})
	}
}
