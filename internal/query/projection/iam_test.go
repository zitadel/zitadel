package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
)

func TestIAMProjection_reduces(t *testing.T) {
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
			name: "reduceGlobalOrgSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.GlobalOrgSetEventType),
					iam.AggregateType,
					[]byte(`{"globalOrgId": "orgid"}`),
				), iam.GlobalOrgSetMapper),
			},
			reduce: (&iamProjection{}).reduceGlobalOrgSet,
			want: wantReduce{
				projection:       IAMProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.iam (id, change_date, sequence, global_org_id) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								uint64(15),
								"orgid",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceGlobalOrgSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.ProjectSetEventType),
					iam.AggregateType,
					[]byte(`{"iamProjectId": "project-id"}`),
				), iam.ProjectSetMapper),
			},
			reduce: (&iamProjection{}).reduceIAMProjectSet,
			want: wantReduce{
				projection:       IAMProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.iam (id, change_date, sequence, iam_project_id) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								uint64(15),
								"project-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSetupStarted",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SetupStartedEventType),
					iam.AggregateType,
					[]byte(`{"Step": 1}`),
				), iam.SetupStepMapper),
			},
			reduce: (&iamProjection{}).reduceSetupEvent,
			want: wantReduce{
				projection:       IAMProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.iam (id, change_date, sequence, setup_started) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								uint64(15),
								domain.Step1,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSetupDone",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SetupDoneEventType),
					iam.AggregateType,
					[]byte(`{"Step": 1}`),
				), iam.SetupStepMapper),
			},
			reduce: (&iamProjection{}).reduceSetupEvent,
			want: wantReduce{
				projection:       IAMProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.iam (id, change_date, sequence, setup_done) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								uint64(15),
								domain.Step1,
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
