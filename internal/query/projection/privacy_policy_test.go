package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestPrivacyPolicyProjection_reduces(t *testing.T) {
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
					repository.EventType(org.PrivacyPolicyAddedEventType),
					org.AggregateType,
					[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link"
}`),
				), org.PrivacyPolicyAddedEventMapper),
			},
			reduce: (&PrivacyPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       PrivacyPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.privacy_policies (creation_date, change_date, sequence, id, state, privacy_link, tos_link, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								"http://privacy.link",
								"http://tos.link",
								false,
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceChanged",
			reduce: (&PrivacyPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.PrivacyPolicyChangedEventType),
					org.AggregateType,
					[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link"
		}`),
				), org.PrivacyPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       PrivacyPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.privacy_policies SET (change_date, sequence, privacy_link, tos_link) = ($1, $2, $3, $4) WHERE (id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"http://privacy.link",
								"http://tos.link",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceRemoved",
			reduce: (&PrivacyPolicyProjection{}).reduceRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.PrivacyPolicyRemovedEventType),
					org.AggregateType,
					nil,
				), org.PrivacyPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       PrivacyPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.privacy_policies WHERE (id = $1)",
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
			reduce: (&PrivacyPolicyProjection{}).reduceAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.PrivacyPolicyAddedEventType),
					instance.AggregateType,
					[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link"
					}`),
				), instance.PrivacyPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       PrivacyPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.privacy_policies (creation_date, change_date, sequence, id, state, privacy_link, tos_link, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								"http://privacy.link",
								"http://tos.link",
								true,
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceChanged",
			reduce: (&PrivacyPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.PrivacyPolicyChangedEventType),
					instance.AggregateType,
					[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link"
					}`),
				), instance.PrivacyPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       PrivacyPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.privacy_policies SET (change_date, sequence, privacy_link, tos_link) = ($1, $2, $3, $4) WHERE (id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"http://privacy.link",
								"http://tos.link",
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
