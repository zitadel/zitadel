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
			name: "org.reduceAdded.Title",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextSetEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, title) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "org.reduceAdded.PreHeader",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextSetEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, pre_header) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "org.reduceAdded.Subject",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextSetEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, subject) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "org.reduceAdded.Greeting",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextSetEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, greeting) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "org.reduceAdded.Text",
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
			reduce: (&messageTextProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "org.reduceAdded.ButtonText",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextSetEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, button_text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "org.reduceAdded.Footer",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextSetEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, footer_text) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "org.reduceRemoved.Title",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextRemovedEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, title) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org.reduceRemoved.PreHeader",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextRemovedEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, pre_header) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org.reduceRemoved.Subject",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextRemovedEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, subject) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org.reduceRemoved.Greeting",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextRemovedEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, greeting) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org.reduceRemoved.Text",
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
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, text) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org.reduceRemoved.ButtonText",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextRemovedEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, button_text) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name: "org.reduceRemoved.Footer",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextRemovedEventType),
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
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, footer_text) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
			name:   "org.reduceRemoved",
			reduce: (&messageTextProjection{}).reduceTemplateRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.CustomTextTemplateRemovedEventType),
					org.AggregateType,
					[]byte(`{
						"key": "Title", 
						"language": "en", 
						"template": "InitCode"
					}`),
				), org.CustomTextTemplateRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.message_texts WHERE (aggregate_id = $1) AND (type = $2) AND (language = $3)",
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
			reduce: (&messageTextProjection{}).reduceAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.CustomTextSetEventType),
					iam.AggregateType,
					[]byte(`{
						"key": "Title",
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
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.message_texts (aggregate_id, creation_date, change_date, sequence, state, type, language, title) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
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
			name: "iam.reduceRemoved.Title",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.CustomTextRemovedEventType),
					iam.AggregateType,
					[]byte(`{
						"key": "Title",
						"language": "en",
						"template": "InitCode"
					}`),
				), iam.CustomTextRemovedEventMapper),
			},
			reduce: (&messageTextProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       MessageTextTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.message_texts SET (change_date, sequence, title) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (type = $5) AND (language = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
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
