package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestUserAuthMethodProjection_reduces(t *testing.T) {
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
			name: "reduceAddedPasswordless",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanPasswordlessTokenAddedType),
					user.AggregateType,
					[]byte(`{
						"webAuthNTokenId": "token-id"
					}`),
				), user.HumanPasswordlessAddedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceInitAuthMethod,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserAuthMethodTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.user_auth_methods (token_id, creation_date, change_date, resource_owner, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"token-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"agg-id",
								uint64(15),
								domain.MFAStateNotReady,
								domain.UserAuthMethodTypePasswordless,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAddedU2F",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanU2FTokenAddedType),
					user.AggregateType,
					[]byte(`{
						"webAuthNTokenId": "token-id"
					}`),
				), user.HumanU2FAddedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceInitAuthMethod,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserAuthMethodTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.user_auth_methods (token_id, creation_date, change_date, resource_owner, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"token-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"agg-id",
								uint64(15),
								domain.MFAStateNotReady,
								domain.UserAuthMethodTypeU2F,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAddedOTP",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanMFAOTPAddedType),
					user.AggregateType,
					[]byte(`{
					}`),
				), user.HumanOTPAddedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceInitAuthMethod,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserAuthMethodTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.user_auth_methods (token_id, creation_date, change_date, resource_owner, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"",
								anyArg{},
								anyArg{},
								"ro-id",
								"agg-id",
								uint64(15),
								domain.MFAStateNotReady,
								domain.UserAuthMethodTypeOTP,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceVerifiedPasswordless",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanPasswordlessTokenVerifiedType),
					user.AggregateType,
					[]byte(`{
						"webAuthNTokenId": "token-id",
						"webAuthNTokenName": "name"
					}`),
				), user.HumanPasswordlessVerifiedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceActivateEvent,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserAuthMethodTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.user_auth_methods SET (change_date, sequence, name, state) = ($1, $2, $3, $4) WHERE (user_id = $5) AND (method_type = $6) AND (resource_owner = $7) AND (token_id = $8)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.MFAStateReady,
								"agg-id",
								domain.UserAuthMethodTypePasswordless,
								"ro-id",
								"token-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceVerifiedU2F",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanU2FTokenVerifiedType),
					user.AggregateType,
					[]byte(`{
						"webAuthNTokenId": "token-id",
						"webAuthNTokenName": "name"
					}`),
				), user.HumanU2FVerifiedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceActivateEvent,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserAuthMethodTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.user_auth_methods SET (change_date, sequence, name, state) = ($1, $2, $3, $4) WHERE (user_id = $5) AND (method_type = $6) AND (resource_owner = $7) AND (token_id = $8)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.MFAStateReady,
								"agg-id",
								domain.UserAuthMethodTypeU2F,
								"ro-id",
								"token-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceVerifiedOTP",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanMFAOTPVerifiedType),
					user.AggregateType,
					[]byte(`{
					}`),
				), user.HumanOTPVerifiedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceActivateEvent,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserAuthMethodTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.user_auth_methods SET (change_date, sequence, name, state) = ($1, $2, $3, $4) WHERE (user_id = $5) AND (method_type = $6) AND (resource_owner = $7) AND (token_id = $8)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
								domain.MFAStateReady,
								"agg-id",
								domain.UserAuthMethodTypeOTP,
								"ro-id",
								"",
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
